package netlink

import (
	"fmt"
	"net"
	"strings"

	"github.com/oss-fun/netlink/nl"
	"golang.org/x/sys/unix"
	
	"github.com/oss-fun/netlink/nlunix"
)

// AddrAdd will add an IP address to a link device.
//
// Equivalent to: `ip addr add $addr dev $link`
//
// If `addr` is an IPv4 address and the broadcast address is not given, it
// will be automatically computed based on the IP mask if /30 or larger.
func AddrAdd(link Link, addr *Addr) error {
	return pkgHandle.AddrAdd(link, addr)
}

// AddrAdd will add an IP address to a link device.
//
// Equivalent to: `ip addr add $addr dev $link`
//
// If `addr` is an IPv4 address and the broadcast address is not given, it
// will be automatically computed based on the IP mask if /30 or larger.
func (h *Handle) AddrAdd(link Link, addr *Addr) error {
	req := h.newNetlinkRequest(unix.RTM_NEWADDR, nlunix.NLM_F_CREATE|nlunix.NLM_F_EXCL|nlunix.NLM_F_ACK)
	return h.addrHandle(link, addr, req)
}

// AddrDel will delete an IP address from a link device.
//
// Equivalent to: `ip addr del $addr dev $link`
//
// If `addr` is an IPv4 address and the broadcast address is not given, it
// will be automatically computed based on the IP mask if /30 or larger.
func AddrDel(link Link, addr *Addr) error {
	return pkgHandle.AddrDel(link, addr)
}

// AddrDel will delete an IP address from a link device.
// Equivalent to: `ip addr del $addr dev $link`
//
// If `addr` is an IPv4 address and the broadcast address is not given, it
// will be automatically computed based on the IP mask if /30 or larger.
func (h *Handle) AddrDel(link Link, addr *Addr) error {
	req := h.newNetlinkRequest(unix.RTM_DELADDR, nlunix.NLM_F_ACK)
	return h.addrHandle(link, addr, req)
}

func (h *Handle) addrHandle(link Link, addr *Addr, req *nl.NetlinkRequest) error {
	family := nl.GetIPFamily(addr.IP)
	msg := nl.NewIfAddrmsg(family)
	msg.Scope = uint8(addr.Scope)
	if link == nil {
		msg.Index = uint32(addr.LinkIndex)
	} else {
		base := link.Attrs()
		if addr.Label != "" && !strings.HasPrefix(addr.Label, base.Name) {
			return fmt.Errorf("label must begin with interface name")
		}
		h.ensureIndex(base)
		msg.Index = uint32(base.Index)
	}
	mask := addr.Mask
	if addr.Peer != nil {
		mask = addr.Peer.Mask
	}
	prefixlen, masklen := mask.Size()
	msg.Prefixlen = uint8(prefixlen)
	req.AddData(msg)

	var localAddrData []byte
	if family == FAMILY_V4 {
		localAddrData = addr.IP.To4()
	} else {
		localAddrData = addr.IP.To16()
	}

	localData := nl.NewRtAttr(nlunix.IFA_LOCAL, localAddrData)
	req.AddData(localData)
	var peerAddrData []byte
	if addr.Peer != nil {
		if family == FAMILY_V4 {
			peerAddrData = addr.Peer.IP.To4()
		} else {
			peerAddrData = addr.Peer.IP.To16()
		}
	} else {
		peerAddrData = localAddrData
	}

	addressData := nl.NewRtAttr(nlunix.IFA_ADDRESS, peerAddrData)
	req.AddData(addressData)

	if addr.Flags != 0 {
		if addr.Flags <= 0xff {
			msg.IfAddrmsg.Flags = uint8(addr.Flags)
		} else {
			b := make([]byte, 4)
			native.PutUint32(b, uint32(addr.Flags))
			flagsData := nl.NewRtAttr(nlunix.IFA_FLAGS, b)
			req.AddData(flagsData)
		}
	}

	if family == FAMILY_V4 {
		// Automatically set the broadcast address if it is unset and the
		// subnet is large enough to sensibly have one (/30 or larger).
		// See: RFC 3021
		if addr.Broadcast == nil && prefixlen < 31 {
			calcBroadcast := make(net.IP, masklen/8)
			for i := range localAddrData {
				calcBroadcast[i] = localAddrData[i] | ^mask[i]
			}
			addr.Broadcast = calcBroadcast
		}

		if addr.Broadcast != nil {
			req.AddData(nl.NewRtAttr(nlunix.IFA_BROADCAST, addr.Broadcast))
		}

		if addr.Label != "" {
			labelData := nl.NewRtAttr(nlunix.IFA_LABEL, nl.ZeroTerminated(addr.Label))
			req.AddData(labelData)
		}
	}

	// 0 is the default value for these attributes. However, 0 means "expired", while the least-surprising default
	// value should be "forever". To compensate for that, only add the attributes if at least one of the values is
	// non-zero, which means the caller has explicitly set them
	if addr.ValidLft > 0 || addr.PreferedLft > 0 {
		cachedata := nl.IfaCacheInfo{nlunix.IfaCacheinfo{
			Valid:    uint32(addr.ValidLft),
			Prefered: uint32(addr.PreferedLft),
		}}
		req.AddData(nl.NewRtAttr(nlunix.IFA_CACHEINFO, cachedata.Serialize()))
	}

	_, err := req.Execute(nlunix.NETLINK_ROUTE, 0)
	return err
}

// AddrList gets a list of IP addresses in the system.
// Equivalent to: `ip addr show`.
// The list can be filtered by link and ip family.
func AddrList(link Link, family int) ([]Addr, error) {
	return pkgHandle.AddrList(link, family)
}

// AddrList gets a list of IP addresses in the system.
// Equivalent to: `ip addr show`.
// The list can be filtered by link and ip family.
func (h *Handle) AddrList(link Link, family int) ([]Addr, error) {
	req := h.newNetlinkRequest(nlunix.RTM_GETADDR, nlunix.NLM_F_DUMP)
	msg := nl.NewIfAddrmsg(family)
	req.AddData(msg)

	msgs, err := req.Execute(nlunix.NETLINK_ROUTE, nlunix.RTM_NEWADDR)
	if err != nil {
		return nil, err
	}

	indexFilter := 0
	if link != nil {
		base := link.Attrs()
		h.ensureIndex(base)
		indexFilter = base.Index
	}

	var res []Addr
	for _, m := range msgs {
		addr, msgFamily, err := parseAddr(m)
		if err != nil {
			return res, err
		}

		if link != nil && addr.LinkIndex != indexFilter {
			// Ignore messages from other interfaces
			continue
		}

		if family != FAMILY_ALL && msgFamily != family {
			continue
		}

		res = append(res, addr)
	}

	return res, nil
}

func parseAddr(m []byte) (addr Addr, family int, err error) {
	msg := nl.DeserializeIfAddrmsg(m)

	family = -1
	addr.LinkIndex = -1

	attrs, err1 := nl.ParseRouteAttr(m[msg.Len():])
	if err1 != nil {
		err = err1
		return
	}

	family = int(msg.Family)
	addr.LinkIndex = int(msg.Index)

	var local, dst *net.IPNet
	for _, attr := range attrs {
		switch attr.Attr.Type {
		case nlunix.IFA_ADDRESS:
			dst = &net.IPNet{
				IP:   attr.Value,
				Mask: net.CIDRMask(int(msg.Prefixlen), 8*len(attr.Value)),
			}
		case nlunix.IFA_LOCAL:
			// iproute2 manual:
			// If a peer address is specified, the local address
			// cannot have a prefix length. The network prefix is
			// associated with the peer rather than with the local
			// address.
			n := 8 * len(attr.Value)
			local = &net.IPNet{
				IP:   attr.Value,
				Mask: net.CIDRMask(n, n),
			}
		case nlunix.IFA_BROADCAST:
			addr.Broadcast = attr.Value
		case nlunix.IFA_LABEL:
			addr.Label = string(attr.Value[:len(attr.Value)-1])
		case nlunix.IFA_FLAGS:
			addr.Flags = int(native.Uint32(attr.Value[0:4]))
		case nlunix.IFA_CACHEINFO:
			ci := nl.DeserializeIfaCacheInfo(attr.Value)
			addr.PreferedLft = int(ci.Prefered)
			addr.ValidLft = int(ci.Valid)
		}
	}

	// libnl addr.c comment:
	// IPv6 sends the local address as IFA_ADDRESS with no
	// IFA_LOCAL, IPv4 sends both IFA_LOCAL and IFA_ADDRESS
	// with IFA_ADDRESS being the peer address if they differ
	//
	// But obviously, as there are IPv6 PtP addresses, too,
	// IFA_LOCAL should also be handled for IPv6.
	if local != nil {
		if family == FAMILY_V4 && dst != nil && local.IP.Equal(dst.IP) {
			addr.IPNet = dst
		} else {
			addr.IPNet = local
			addr.Peer = dst
		}
	} else {
		addr.IPNet = dst
	}

	addr.Scope = int(msg.Scope)

	return
}

