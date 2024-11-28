package netlink

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
	"syscall"

	"github.com/oss-fun/netlink/nl"
	"github.com/vishvananda/netns"
	"golang.org/x/sys/unix"

	"github.com/oss-fun/netlink/nlunix"
)

// RtAttr is shared so it is in netlink_linux.go

const (
	SCOPE_UNIVERSE Scope = nlunix.RT_SCOPE_UNIVERSE
	SCOPE_SITE     Scope = nlunix.RT_SCOPE_SITE
	SCOPE_LINK     Scope = nlunix.RT_SCOPE_LINK
	SCOPE_HOST     Scope = nlunix.RT_SCOPE_HOST
	SCOPE_NOWHERE  Scope = nlunix.RT_SCOPE_NOWHERE
)

func (s Scope) String() string {
	switch s {
	case SCOPE_UNIVERSE:
		return "universe"
	case SCOPE_SITE:
		return "site"
	case SCOPE_LINK:
		return "link"
	case SCOPE_HOST:
		return "host"
	case SCOPE_NOWHERE:
		return "nowhere"
	default:
		return "unknown"
	}
}

const (
	FLAG_ONLINK    NextHopFlag = nlunix.RTNH_F_ONLINK
	FLAG_PERVASIVE NextHopFlag = nlunix.RTNH_F_PERVASIVE
)

var testFlags = []flagString{
	{f: FLAG_ONLINK, s: "onlink"},
	{f: FLAG_PERVASIVE, s: "pervasive"},
}

func listFlags(flag int) []string {
	var flags []string
	for _, tf := range testFlags {
		if flag&int(tf.f) != 0 {
			flags = append(flags, tf.s)
		}
	}
	return flags
}

func (r *Route) ListFlags() []string {
	return listFlags(r.Flags)
}

func (n *NexthopInfo) ListFlags() []string {
	return listFlags(n.Flags)
}

type MPLSDestination struct {
	Labels []int
}

func (d *MPLSDestination) Family() int {
	return nl.FAMILY_MPLS
}

func (d *MPLSDestination) Decode(buf []byte) error {
	d.Labels = nl.DecodeMPLSStack(buf)
	return nil
}

func (d *MPLSDestination) Encode() ([]byte, error) {
	return nl.EncodeMPLSStack(d.Labels...), nil
}

func (d *MPLSDestination) String() string {
	s := make([]string, 0, len(d.Labels))
	for _, l := range d.Labels {
		s = append(s, fmt.Sprintf("%d", l))
	}
	return strings.Join(s, "/")
}

func (d *MPLSDestination) Equal(x Destination) bool {
	o, ok := x.(*MPLSDestination)
	if !ok {
		return false
	}
	if d == nil && o == nil {
		return true
	}
	if d == nil || o == nil {
		return false
	}
	if d.Labels == nil && o.Labels == nil {
		return true
	}
	if d.Labels == nil || o.Labels == nil {
		return false
	}
	if len(d.Labels) != len(o.Labels) {
		return false
	}
	for i := range d.Labels {
		if d.Labels[i] != o.Labels[i] {
			return false
		}
	}
	return true
}

type Via struct {
	AddrFamily int
	Addr       net.IP
}

func (v *Via) Equal(x Destination) bool {
	o, ok := x.(*Via)
	if !ok {
		return false
	}
	if v.AddrFamily == x.Family() && v.Addr.Equal(o.Addr) {
		return true
	}
	return false
}

func (v *Via) String() string {
	return fmt.Sprintf("Family: %d, Address: %s", v.AddrFamily, v.Addr.String())
}

func (v *Via) Family() int {
	return v.AddrFamily
}

func (v *Via) Encode() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, native, uint16(v.AddrFamily))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, native, v.Addr)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (v *Via) Decode(b []byte) error {
	if len(b) < 6 {
		return fmt.Errorf("decoding failed: buffer too small (%d bytes)", len(b))
	}
	v.AddrFamily = int(native.Uint16(b[0:2]))
	if v.AddrFamily == nl.FAMILY_V4 {
		v.Addr = net.IP(b[2:6])
		return nil
	} else if v.AddrFamily == nl.FAMILY_V6 {
		if len(b) < 18 {
			return fmt.Errorf("decoding failed: buffer too small (%d bytes)", len(b))
		}
		v.Addr = net.IP(b[2:])
		return nil
	}
	return fmt.Errorf("decoding failed: address family %d unknown", v.AddrFamily)
}

// RouteAdd will add a route to the system.
// Equivalent to: `ip route add $route`
func RouteAdd(route *Route) error {
	return pkgHandle.RouteAdd(route)
}

// RouteAdd will add a route to the system.
// Equivalent to: `ip route add $route`
func (h *Handle) RouteAdd(route *Route) error {
	flags := nlunix.NLM_F_CREATE | nlunix.NLM_F_EXCL | nlunix.NLM_F_ACK
	req := h.newNetlinkRequest(nlunix.RTM_NEWROUTE, flags)
	_, err := h.routeHandle(route, req, nl.NewRtMsg())
	return err
}

// RouteAppend will append a route to the system.
// Equivalent to: `ip route append $route`
func RouteAppend(route *Route) error {
	return pkgHandle.RouteAppend(route)
}

// RouteAppend will append a route to the system.
// Equivalent to: `ip route append $route`
func (h *Handle) RouteAppend(route *Route) error {
	flags := nlunix.NLM_F_CREATE | nlunix.NLM_F_APPEND | nlunix.NLM_F_ACK
	req := h.newNetlinkRequest(nlunix.RTM_NEWROUTE, flags)
	_, err := h.routeHandle(route, req, nl.NewRtMsg())
	return err
}

// RouteAddEcmp will add a route to the system.
func RouteAddEcmp(route *Route) error {
	return pkgHandle.RouteAddEcmp(route)
}

// RouteAddEcmp will add a route to the system.
func (h *Handle) RouteAddEcmp(route *Route) error {
	flags := nlunix.NLM_F_CREATE | nlunix.NLM_F_ACK
	req := h.newNetlinkRequest(nlunix.RTM_NEWROUTE, flags)
	_, err := h.routeHandle(route, req, nl.NewRtMsg())
	return err
}

// RouteChange will change an existing route in the system.
// Equivalent to: `ip route change $route`
func RouteChange(route *Route) error {
	return pkgHandle.RouteChange(route)
}

// RouteChange will change an existing route in the system.
// Equivalent to: `ip route change $route`
func (h *Handle) RouteChange(route *Route) error {
	flags := nlunix.NLM_F_REPLACE | nlunix.NLM_F_ACK
	req := h.newNetlinkRequest(nlunix.RTM_NEWROUTE, flags)
	_, err := h.routeHandle(route, req, nl.NewRtMsg())
	return err
}

// RouteReplace will add a route to the system.
// Equivalent to: `ip route replace $route`
func RouteReplace(route *Route) error {
	return pkgHandle.RouteReplace(route)
}

// RouteReplace will add a route to the system.
// Equivalent to: `ip route replace $route`
func (h *Handle) RouteReplace(route *Route) error {
	flags := nlunix.NLM_F_CREATE | nlunix.NLM_F_REPLACE | nlunix.NLM_F_ACK
	req := h.newNetlinkRequest(nlunix.RTM_NEWROUTE, flags)
	_, err := h.routeHandle(route, req, nl.NewRtMsg())
	return err
}

// RouteDel will delete a route from the system.
// Equivalent to: `ip route del $route`
func RouteDel(route *Route) error {
	return pkgHandle.RouteDel(route)
}

// RouteDel will delete a route from the system.
// Equivalent to: `ip route del $route`
func (h *Handle) RouteDel(route *Route) error {
	req := h.newNetlinkRequest(nlunix.RTM_DELROUTE, nlunix.NLM_F_ACK)
	_, err := h.routeHandle(route, req, nl.NewRtDelMsg())
	return err
}

func (h *Handle) routeHandle(route *Route, req *nl.NetlinkRequest, msg *nl.RtMsg) ([][]byte, error) {
	if err := h.prepareRouteReq(route, req, msg); err != nil {
		return nil, err
	}
	return req.Execute(nlunix.NETLINK_ROUTE, 0)
}

func (h *Handle) routeHandleIter(route *Route, req *nl.NetlinkRequest, msg *nl.RtMsg, f func(msg []byte) bool) error {
	if err := h.prepareRouteReq(route, req, msg); err != nil {
		return err
	}
	return req.ExecuteIter(nlunix.NETLINK_ROUTE, 0, f)
}

func (h *Handle) prepareRouteReq(route *Route, req *nl.NetlinkRequest, msg *nl.RtMsg) error {
	if req.NlMsghdr.Type != nlunix.RTM_GETROUTE && (route.Dst == nil || route.Dst.IP == nil) && route.Src == nil && route.Gw == nil && route.MPLSDst == nil {
		return fmt.Errorf("either Dst.IP, Src.IP or Gw must be set")
	}

	family := -1
	var rtAttrs []*nl.RtAttr

	if route.Dst != nil && route.Dst.IP != nil {
		dstLen, _ := route.Dst.Mask.Size()
		msg.Dst_len = uint8(dstLen)
		dstFamily := nl.GetIPFamily(route.Dst.IP)
		family = dstFamily
		var dstData []byte
		if dstFamily == FAMILY_V4 {
			dstData = route.Dst.IP.To4()
		} else {
			dstData = route.Dst.IP.To16()
		}
		rtAttrs = append(rtAttrs, nl.NewRtAttr(unix.RTA_DST, dstData))
	} else if route.MPLSDst != nil {
		family = nl.FAMILY_MPLS
		msg.Dst_len = uint8(20)
		msg.Type = nlunix.RTN_UNICAST
		rtAttrs = append(rtAttrs, nl.NewRtAttr(unix.RTA_DST, nl.EncodeMPLSStack(*route.MPLSDst)))
	}

	if route.NewDst != nil {
		if family != -1 && family != route.NewDst.Family() {
			return fmt.Errorf("new destination and destination are not the same address family")
		}
		buf, err := route.NewDst.Encode()
		if err != nil {
			return err
		}
		rtAttrs = append(rtAttrs, nl.NewRtAttr(nlunix.RTA_NEWDST, buf))
	}

	if route.Encap != nil {
		buf := make([]byte, 2)
		native.PutUint16(buf, uint16(route.Encap.Type()))
		rtAttrs = append(rtAttrs, nl.NewRtAttr(nlunix.RTA_ENCAP_TYPE, buf))
		buf, err := route.Encap.Encode()
		if err != nil {
			return err
		}
		switch route.Encap.Type() {
		case nl.LWTUNNEL_ENCAP_BPF:
			rtAttrs = append(rtAttrs, nl.NewRtAttr(nlunix.RTA_ENCAP|nlunix.NLA_F_NESTED, buf))
		default:
			rtAttrs = append(rtAttrs, nl.NewRtAttr(nlunix.RTA_ENCAP, buf))
		}

	}

	if route.Src != nil {
		srcFamily := nl.GetIPFamily(route.Src)
		if family != -1 && family != srcFamily {
			return fmt.Errorf("source and destination ip are not the same IP family")
		}
		family = srcFamily
		var srcData []byte
		if srcFamily == FAMILY_V4 {
			srcData = route.Src.To4()
		} else {
			srcData = route.Src.To16()
		}
		// The commonly used src ip for routes is actually PREFSRC
		rtAttrs = append(rtAttrs, nl.NewRtAttr(nlunix.RTA_PREFSRC, srcData))
	}

	if route.Gw != nil {
		gwFamily := nl.GetIPFamily(route.Gw)
		if family != -1 && family != gwFamily {
			return fmt.Errorf("gateway, source, and destination ip are not the same IP family")
		}
		family = gwFamily
		var gwData []byte
		if gwFamily == FAMILY_V4 {
			gwData = route.Gw.To4()
		} else {
			gwData = route.Gw.To16()
		}
		rtAttrs = append(rtAttrs, nl.NewRtAttr(unix.RTA_GATEWAY, gwData))
	}

	if route.Via != nil {
		buf, err := route.Via.Encode()
		if err != nil {
			return fmt.Errorf("failed to encode RTA_VIA: %v", err)
		}
		rtAttrs = append(rtAttrs, nl.NewRtAttr(nlunix.RTA_VIA, buf))
	}

	if len(route.MultiPath) > 0 {
		buf := []byte{}
		for _, nh := range route.MultiPath {
			rtnh := &nl.RtNexthop{
				RtNexthop: nlunix.RtNexthop{
					Hops:    uint8(nh.Hops),
					Ifindex: int32(nh.LinkIndex),
					Flags:   uint8(nh.Flags),
				},
			}
			children := []nl.NetlinkRequestData{}
			if nh.Gw != nil {
				gwFamily := nl.GetIPFamily(nh.Gw)
				if family != -1 && family != gwFamily {
					return fmt.Errorf("gateway, source, and destination ip are not the same IP family")
				}
				if gwFamily == FAMILY_V4 {
					children = append(children, nl.NewRtAttr(unix.RTA_GATEWAY, []byte(nh.Gw.To4())))
				} else {
					children = append(children, nl.NewRtAttr(unix.RTA_GATEWAY, []byte(nh.Gw.To16())))
				}
			}
			if nh.NewDst != nil {
				if family != -1 && family != nh.NewDst.Family() {
					return fmt.Errorf("new destination and destination are not the same address family")
				}
				buf, err := nh.NewDst.Encode()
				if err != nil {
					return err
				}
				children = append(children, nl.NewRtAttr(nlunix.RTA_NEWDST, buf))
			}
			if nh.Encap != nil {
				buf := make([]byte, 2)
				native.PutUint16(buf, uint16(nh.Encap.Type()))
				children = append(children, nl.NewRtAttr(nlunix.RTA_ENCAP_TYPE, buf))
				buf, err := nh.Encap.Encode()
				if err != nil {
					return err
				}
				children = append(children, nl.NewRtAttr(nlunix.RTA_ENCAP, buf))
			}
			if nh.Via != nil {
				buf, err := nh.Via.Encode()
				if err != nil {
					return err
				}
				children = append(children, nl.NewRtAttr(nlunix.RTA_VIA, buf))
			}
			rtnh.Children = children
			buf = append(buf, rtnh.Serialize()...)
		}
		rtAttrs = append(rtAttrs, nl.NewRtAttr(nlunix.RTA_MULTIPATH, buf))
	}

	if route.Table > 0 {
		if route.Table >= 256 {
			msg.Table = nlunix.RT_TABLE_UNSPEC
			b := make([]byte, 4)
			native.PutUint32(b, uint32(route.Table))
			rtAttrs = append(rtAttrs, nl.NewRtAttr(nlunix.RTA_TABLE, b))
		} else {
			msg.Table = uint8(route.Table)
		}
	}

	if route.Priority > 0 {
		b := make([]byte, 4)
		native.PutUint32(b, uint32(route.Priority))
		rtAttrs = append(rtAttrs, nl.NewRtAttr(nlunix.RTA_PRIORITY, b))
	}
	if route.Realm > 0 {
		b := make([]byte, 4)
		native.PutUint32(b, uint32(route.Realm))
		rtAttrs = append(rtAttrs, nl.NewRtAttr(nlunix.RTA_FLOW, b))
	}
	if route.Tos > 0 {
		msg.Tos = uint8(route.Tos)
	}
	if route.Protocol > 0 {
		msg.Protocol = uint8(route.Protocol)
	}
	if route.Type > 0 {
		msg.Type = uint8(route.Type)
	}

	var metrics []*nl.RtAttr
	if route.MTU > 0 {
		b := nl.Uint32Attr(uint32(route.MTU))
		metrics = append(metrics, nl.NewRtAttr(nlunix.RTAX_MTU, b))
	}
	if route.Window > 0 {
		b := nl.Uint32Attr(uint32(route.Window))
		metrics = append(metrics, nl.NewRtAttr(nlunix.RTAX_WINDOW, b))
	}
	if route.Rtt > 0 {
		b := nl.Uint32Attr(uint32(route.Rtt))
		metrics = append(metrics, nl.NewRtAttr(nlunix.RTAX_RTT, b))
	}
	if route.RttVar > 0 {
		b := nl.Uint32Attr(uint32(route.RttVar))
		metrics = append(metrics, nl.NewRtAttr(nlunix.RTAX_RTTVAR, b))
	}
	if route.Ssthresh > 0 {
		b := nl.Uint32Attr(uint32(route.Ssthresh))
		metrics = append(metrics, nl.NewRtAttr(nlunix.RTAX_SSTHRESH, b))
	}
	if route.Cwnd > 0 {
		b := nl.Uint32Attr(uint32(route.Cwnd))
		metrics = append(metrics, nl.NewRtAttr(nlunix.RTAX_CWND, b))
	}
	if route.AdvMSS > 0 {
		b := nl.Uint32Attr(uint32(route.AdvMSS))
		metrics = append(metrics, nl.NewRtAttr(nlunix.RTAX_ADVMSS, b))
	}
	if route.Reordering > 0 {
		b := nl.Uint32Attr(uint32(route.Reordering))
		metrics = append(metrics, nl.NewRtAttr(unix.RTAX_REORDERING, b))
	}
	if route.Hoplimit > 0 {
		b := nl.Uint32Attr(uint32(route.Hoplimit))
		metrics = append(metrics, nl.NewRtAttr(unix.RTAX_HOPLIMIT, b))
	}
	if route.InitCwnd > 0 {
		b := nl.Uint32Attr(uint32(route.InitCwnd))
		metrics = append(metrics, nl.NewRtAttr(unix.RTAX_INITCWND, b))
	}
	if route.Features > 0 {
		b := nl.Uint32Attr(uint32(route.Features))
		metrics = append(metrics, nl.NewRtAttr(unix.RTAX_FEATURES, b))
	}
	if route.RtoMin > 0 {
		b := nl.Uint32Attr(uint32(route.RtoMin))
		metrics = append(metrics, nl.NewRtAttr(unix.RTAX_RTO_MIN, b))
	}
	if route.InitRwnd > 0 {
		b := nl.Uint32Attr(uint32(route.InitRwnd))
		metrics = append(metrics, nl.NewRtAttr(unix.RTAX_INITRWND, b))
	}
	if route.QuickACK > 0 {
		b := nl.Uint32Attr(uint32(route.QuickACK))
		metrics = append(metrics, nl.NewRtAttr(unix.RTAX_QUICKACK, b))
	}
	if route.Congctl != "" {
		b := nl.ZeroTerminated(route.Congctl)
		metrics = append(metrics, nl.NewRtAttr(unix.RTAX_CC_ALGO, b))
	}
	if route.FastOpenNoCookie > 0 {
		b := nl.Uint32Attr(uint32(route.FastOpenNoCookie))
		metrics = append(metrics, nl.NewRtAttr(unix.RTAX_FASTOPEN_NO_COOKIE, b))
	}

	if metrics != nil {
		attr := nl.NewRtAttr(unix.RTA_METRICS, nil)
		for _, metric := range metrics {
			attr.AddChild(metric)
		}
		rtAttrs = append(rtAttrs, attr)
	}

	msg.Flags = uint32(route.Flags)
	msg.Scope = uint8(route.Scope)
	// only overwrite family if it was not set in msg
	if msg.Family == 0 {
		msg.Family = uint8(family)
	}
	req.AddData(msg)
	for _, attr := range rtAttrs {
		req.AddData(attr)
	}

	if (req.NlMsghdr.Type != unix.RTM_GETROUTE) || (req.NlMsghdr.Type == unix.RTM_GETROUTE && route.LinkIndex > 0) {
		b := make([]byte, 4)
		native.PutUint32(b, uint32(route.LinkIndex))
		req.AddData(nl.NewRtAttr(unix.RTA_OIF, b))
	}
	return nil
}

// RouteList gets a list of routes in the system.
// Equivalent to: `ip route show`.
// The list can be filtered by link and ip family.
func RouteList(link Link, family int) ([]Route, error) {
	return pkgHandle.RouteList(link, family)
}

// RouteList gets a list of routes in the system.
// Equivalent to: `ip route show`.
// The list can be filtered by link and ip family.
func (h *Handle) RouteList(link Link, family int) ([]Route, error) {
	routeFilter := &Route{}
	if link != nil {
		routeFilter.LinkIndex = link.Attrs().Index

		return h.RouteListFiltered(family, routeFilter, RT_FILTER_OIF)
	}
	return h.RouteListFiltered(family, routeFilter, 0)
}

// RouteListFiltered gets a list of routes in the system filtered with specified rules.
// All rules must be defined in RouteFilter struct
func RouteListFiltered(family int, filter *Route, filterMask uint64) ([]Route, error) {
	return pkgHandle.RouteListFiltered(family, filter, filterMask)
}

// RouteListFiltered gets a list of routes in the system filtered with specified rules.
// All rules must be defined in RouteFilter struct
func (h *Handle) RouteListFiltered(family int, filter *Route, filterMask uint64) ([]Route, error) {
	var res []Route
	err := h.RouteListFilteredIter(family, filter, filterMask, func(route Route) (cont bool) {
		res = append(res, route)
		return true
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// RouteListFilteredIter passes each route that matches the filter to the given iterator func.  Iteration continues
// until all routes are loaded or the func returns false.
func RouteListFilteredIter(family int, filter *Route, filterMask uint64, f func(Route) (cont bool)) error {
	return pkgHandle.RouteListFilteredIter(family, filter, filterMask, f)
}

func (h *Handle) RouteListFilteredIter(family int, filter *Route, filterMask uint64, f func(Route) (cont bool)) error {
	req := h.newNetlinkRequest(unix.RTM_GETROUTE, unix.NLM_F_DUMP)
	rtmsg := &nl.RtMsg{}
	rtmsg.Family = uint8(family)

	var parseErr error
	err := h.routeHandleIter(filter, req, rtmsg, func(m []byte) bool {
		msg := nl.DeserializeRtMsg(m)
		if family != FAMILY_ALL && msg.Family != uint8(family) {
			// Ignore routes not matching requested family
			return true
		}
		if msg.Flags&unix.RTM_F_CLONED != 0 {
			// Ignore cloned routes
			return true
		}
		if msg.Table != unix.RT_TABLE_MAIN {
			if filter == nil || filterMask&RT_FILTER_TABLE == 0 {
				// Ignore non-main tables
				return true
			}
		}
		route, err := deserializeRoute(m)
		if err != nil {
			parseErr = err
			return false
		}
		if filter != nil {
			switch {
			case filterMask&RT_FILTER_TABLE != 0 && filter.Table != unix.RT_TABLE_UNSPEC && route.Table != filter.Table:
				return true
			case filterMask&RT_FILTER_PROTOCOL != 0 && route.Protocol != filter.Protocol:
				return true
			case filterMask&RT_FILTER_SCOPE != 0 && route.Scope != filter.Scope:
				return true
			case filterMask&RT_FILTER_TYPE != 0 && route.Type != filter.Type:
				return true
			case filterMask&RT_FILTER_TOS != 0 && route.Tos != filter.Tos:
				return true
			case filterMask&RT_FILTER_REALM != 0 && route.Realm != filter.Realm:
				return true
			case filterMask&RT_FILTER_OIF != 0 && route.LinkIndex != filter.LinkIndex:
				return true
			case filterMask&RT_FILTER_IIF != 0 && route.ILinkIndex != filter.ILinkIndex:
				return true
			case filterMask&RT_FILTER_GW != 0 && !route.Gw.Equal(filter.Gw):
				return true
			case filterMask&RT_FILTER_SRC != 0 && !route.Src.Equal(filter.Src):
				return true
			case filterMask&RT_FILTER_DST != 0:
				if filter.MPLSDst == nil || route.MPLSDst == nil || (*filter.MPLSDst) != (*route.MPLSDst) {
					if filter.Dst == nil {
						filter.Dst = genZeroIPNet(family)
					}
					if !ipNetEqual(route.Dst, filter.Dst) {
						return true
					}
				}
			case filterMask&RT_FILTER_HOPLIMIT != 0 && route.Hoplimit != filter.Hoplimit:
				return true
			}
		}
		return f(route)
	})
	if err != nil {
		return err
	}
	if parseErr != nil {
		return parseErr
	}
	return nil
}

// deserializeRoute decodes a binary netlink message into a Route struct
func deserializeRoute(m []byte) (Route, error) {
	msg := nl.DeserializeRtMsg(m)
	attrs, err := nl.ParseRouteAttr(m[msg.Len():])
	if err != nil {
		return Route{}, err
	}
	route := Route{
		Scope:    Scope(msg.Scope),
		Protocol: RouteProtocol(int(msg.Protocol)),
		Table:    int(msg.Table),
		Type:     int(msg.Type),
		Tos:      int(msg.Tos),
		Flags:    int(msg.Flags),
		Family:   int(msg.Family),
	}

	var encap, encapType syscall.NetlinkRouteAttr
	for _, attr := range attrs {
		switch attr.Attr.Type {
		case unix.RTA_GATEWAY:
			route.Gw = net.IP(attr.Value)
		case unix.RTA_PREFSRC:
			route.Src = net.IP(attr.Value)
		case unix.RTA_DST:
			if msg.Family == nl.FAMILY_MPLS {
				stack := nl.DecodeMPLSStack(attr.Value)
				if len(stack) == 0 || len(stack) > 1 {
					return route, fmt.Errorf("invalid MPLS RTA_DST")
				}
				route.MPLSDst = &stack[0]
			} else {
				route.Dst = &net.IPNet{
					IP:   attr.Value,
					Mask: net.CIDRMask(int(msg.Dst_len), 8*len(attr.Value)),
				}
			}
		case unix.RTA_OIF:
			route.LinkIndex = int(native.Uint32(attr.Value[0:4]))
		case unix.RTA_IIF:
			route.ILinkIndex = int(native.Uint32(attr.Value[0:4]))
		case unix.RTA_PRIORITY:
			route.Priority = int(native.Uint32(attr.Value[0:4]))
		case unix.RTA_FLOW:
			route.Realm = int(native.Uint32(attr.Value[0:4]))
		case unix.RTA_TABLE:
			route.Table = int(native.Uint32(attr.Value[0:4]))
		case unix.RTA_MULTIPATH:
			parseRtNexthop := func(value []byte) (*NexthopInfo, []byte, error) {
				if len(value) < unix.SizeofRtNexthop {
					return nil, nil, fmt.Errorf("lack of bytes")
				}
				nh := nl.DeserializeRtNexthop(value)
				if len(value) < int(nh.RtNexthop.Len) {
					return nil, nil, fmt.Errorf("lack of bytes")
				}
				info := &NexthopInfo{
					LinkIndex: int(nh.RtNexthop.Ifindex),
					Hops:      int(nh.RtNexthop.Hops),
					Flags:     int(nh.RtNexthop.Flags),
				}
				attrs, err := nl.ParseRouteAttr(value[unix.SizeofRtNexthop:int(nh.RtNexthop.Len)])
				if err != nil {
					return nil, nil, err
				}
				var encap, encapType syscall.NetlinkRouteAttr
				for _, attr := range attrs {
					switch attr.Attr.Type {
					case unix.RTA_GATEWAY:
						info.Gw = net.IP(attr.Value)
					case unix.RTA_NEWDST:
						var d Destination
						switch msg.Family {
						case nl.FAMILY_MPLS:
							d = &MPLSDestination{}
						}
						if err := d.Decode(attr.Value); err != nil {
							return nil, nil, err
						}
						info.NewDst = d
					case unix.RTA_ENCAP_TYPE:
						encapType = attr
					case unix.RTA_ENCAP:
						encap = attr
					case unix.RTA_VIA:
						d := &Via{}
						if err := d.Decode(attr.Value); err != nil {
							return nil, nil, err
						}
						info.Via = d
					}
				}

				if len(encap.Value) != 0 && len(encapType.Value) != 0 {
					typ := int(native.Uint16(encapType.Value[0:2]))
					var e Encap
					switch typ {
					case nl.LWTUNNEL_ENCAP_MPLS:
						e = &MPLSEncap{}
						if err := e.Decode(encap.Value); err != nil {
							return nil, nil, err
						}
					}
					info.Encap = e
				}

				return info, value[int(nh.RtNexthop.Len):], nil
			}
			rest := attr.Value
			for len(rest) > 0 {
				info, buf, err := parseRtNexthop(rest)
				if err != nil {
					return route, err
				}
				route.MultiPath = append(route.MultiPath, info)
				rest = buf
			}
		case unix.RTA_NEWDST:
			var d Destination
			switch msg.Family {
			case nl.FAMILY_MPLS:
				d = &MPLSDestination{}
			}
			if err := d.Decode(attr.Value); err != nil {
				return route, err
			}
			route.NewDst = d
		case unix.RTA_VIA:
			v := &Via{}
			if err := v.Decode(attr.Value); err != nil {
				return route, err
			}
			route.Via = v
		case unix.RTA_ENCAP_TYPE:
			encapType = attr
		case unix.RTA_ENCAP:
			encap = attr
		case unix.RTA_METRICS:
			metrics, err := nl.ParseRouteAttr(attr.Value)
			if err != nil {
				return route, err
			}
			for _, metric := range metrics {
				switch metric.Attr.Type {
				case unix.RTAX_MTU:
					route.MTU = int(native.Uint32(metric.Value[0:4]))
				case unix.RTAX_WINDOW:
					route.Window = int(native.Uint32(metric.Value[0:4]))
				case unix.RTAX_RTT:
					route.Rtt = int(native.Uint32(metric.Value[0:4]))
				case unix.RTAX_RTTVAR:
					route.RttVar = int(native.Uint32(metric.Value[0:4]))
				case unix.RTAX_SSTHRESH:
					route.Ssthresh = int(native.Uint32(metric.Value[0:4]))
				case unix.RTAX_CWND:
					route.Cwnd = int(native.Uint32(metric.Value[0:4]))
				case unix.RTAX_ADVMSS:
					route.AdvMSS = int(native.Uint32(metric.Value[0:4]))
				case unix.RTAX_REORDERING:
					route.Reordering = int(native.Uint32(metric.Value[0:4]))
				case unix.RTAX_HOPLIMIT:
					route.Hoplimit = int(native.Uint32(metric.Value[0:4]))
				case unix.RTAX_INITCWND:
					route.InitCwnd = int(native.Uint32(metric.Value[0:4]))
				case unix.RTAX_FEATURES:
					route.Features = int(native.Uint32(metric.Value[0:4]))
				case unix.RTAX_RTO_MIN:
					route.RtoMin = int(native.Uint32(metric.Value[0:4]))
				case unix.RTAX_INITRWND:
					route.InitRwnd = int(native.Uint32(metric.Value[0:4]))
				case unix.RTAX_QUICKACK:
					route.QuickACK = int(native.Uint32(metric.Value[0:4]))
				case unix.RTAX_CC_ALGO:
					route.Congctl = nl.BytesToString(metric.Value)
				case unix.RTAX_FASTOPEN_NO_COOKIE:
					route.FastOpenNoCookie = int(native.Uint32(metric.Value[0:4]))
				}
			}
		}
	}

	// Same logic to generate "default" dst with iproute2 implementation
	if route.Dst == nil {
		var addLen int
		var ip net.IP
		switch msg.Family {
		case FAMILY_V4:
			addLen = net.IPv4len
			ip = net.IPv4zero
		case FAMILY_V6:
			addLen = net.IPv6len
			ip = net.IPv6zero
		}

		if addLen != 0 {
			route.Dst = &net.IPNet{
				IP:   ip,
				Mask: net.CIDRMask(int(msg.Dst_len), 8*addLen),
			}
		}
	}

	if len(encap.Value) != 0 && len(encapType.Value) != 0 {
		typ := int(native.Uint16(encapType.Value[0:2]))
		var e Encap
		switch typ {
		case nl.LWTUNNEL_ENCAP_MPLS:
			e = &MPLSEncap{}
			if err := e.Decode(encap.Value); err != nil {
				return route, err
			}
		case nl.LWTUNNEL_ENCAP_SEG6:
			e = &SEG6Encap{}
			if err := e.Decode(encap.Value); err != nil {
				return route, err
			}
		case nl.LWTUNNEL_ENCAP_SEG6_LOCAL:
			e = &SEG6LocalEncap{}
			if err := e.Decode(encap.Value); err != nil {
				return route, err
			}
		case nl.LWTUNNEL_ENCAP_BPF:
			e = &BpfEncap{}
			if err := e.Decode(encap.Value); err != nil {
				return route, err
			}
		}
		route.Encap = e
	}

	return route, nil
}

// RouteGetOptions contains a set of options to use with
// RouteGetWithOptions
type RouteGetOptions struct {
	Iif      string
	IifIndex int
	Oif      string
	VrfName  string
	SrcAddr  net.IP
	UID      *uint32
	Mark     uint32
	FIBMatch bool
}

// RouteGetWithOptions gets a route to a specific destination from the host system.
// Equivalent to: 'ip route get <> vrf <VrfName>'.
func RouteGetWithOptions(destination net.IP, options *RouteGetOptions) ([]Route, error) {
	return pkgHandle.RouteGetWithOptions(destination, options)
}

// RouteGet gets a route to a specific destination from the host system.
// Equivalent to: 'ip route get'.
func RouteGet(destination net.IP) ([]Route, error) {
	return pkgHandle.RouteGet(destination)
}

// RouteGetWithOptions gets a route to a specific destination from the host system.
// Equivalent to: 'ip route get <> vrf <VrfName>'.
func (h *Handle) RouteGetWithOptions(destination net.IP, options *RouteGetOptions) ([]Route, error) {
	req := h.newNetlinkRequest(unix.RTM_GETROUTE, unix.NLM_F_REQUEST)
	family := nl.GetIPFamily(destination)
	var destinationData []byte
	var bitlen uint8
	if family == FAMILY_V4 {
		destinationData = destination.To4()
		bitlen = 32
	} else {
		destinationData = destination.To16()
		bitlen = 128
	}
	msg := &nl.RtMsg{}
	msg.Family = uint8(family)
	msg.Dst_len = bitlen
	if options != nil && options.SrcAddr != nil {
		msg.Src_len = bitlen
	}
	msg.Flags = unix.RTM_F_LOOKUP_TABLE
	if options != nil && options.FIBMatch {
		msg.Flags |= unix.RTM_F_FIB_MATCH
	}
	req.AddData(msg)

	rtaDst := nl.NewRtAttr(unix.RTA_DST, destinationData)
	req.AddData(rtaDst)

	if options != nil {
		if options.VrfName != "" {
			link, err := h.LinkByName(options.VrfName)
			if err != nil {
				return nil, err
			}
			b := make([]byte, 4)
			native.PutUint32(b, uint32(link.Attrs().Index))

			req.AddData(nl.NewRtAttr(unix.RTA_OIF, b))
		}

		iifIndex := 0
		if len(options.Iif) > 0 {
			link, err := h.LinkByName(options.Iif)
			if err != nil {
				return nil, err
			}

			iifIndex = link.Attrs().Index
		} else if options.IifIndex > 0 {
			iifIndex = options.IifIndex
		}

		if iifIndex > 0 {
			b := make([]byte, 4)
			native.PutUint32(b, uint32(iifIndex))

			req.AddData(nl.NewRtAttr(unix.RTA_IIF, b))
		}

		if len(options.Oif) > 0 {
			link, err := h.LinkByName(options.Oif)
			if err != nil {
				return nil, err
			}

			b := make([]byte, 4)
			native.PutUint32(b, uint32(link.Attrs().Index))

			req.AddData(nl.NewRtAttr(unix.RTA_OIF, b))
		}

		if options.SrcAddr != nil {
			var srcAddr []byte
			if family == FAMILY_V4 {
				srcAddr = options.SrcAddr.To4()
			} else {
				srcAddr = options.SrcAddr.To16()
			}

			req.AddData(nl.NewRtAttr(unix.RTA_SRC, srcAddr))
		}

		if options.UID != nil {
			uid := *options.UID
			b := make([]byte, 4)
			native.PutUint32(b, uid)

			req.AddData(nl.NewRtAttr(unix.RTA_UID, b))
		}

		if options.Mark > 0 {
			b := make([]byte, 4)
			native.PutUint32(b, options.Mark)

			req.AddData(nl.NewRtAttr(unix.RTA_MARK, b))
		}
	}

	msgs, err := req.Execute(unix.NETLINK_ROUTE, unix.RTM_NEWROUTE)
	if err != nil {
		return nil, err
	}

	var res []Route
	for _, m := range msgs {
		route, err := deserializeRoute(m)
		if err != nil {
			return nil, err
		}
		res = append(res, route)
	}
	return res, nil
}

// RouteGet gets a route to a specific destination from the host system.
// Equivalent to: 'ip route get'.
func (h *Handle) RouteGet(destination net.IP) ([]Route, error) {
	return h.RouteGetWithOptions(destination, nil)
}

// RouteSubscribe takes a chan down which notifications will be sent
// when routes are added or deleted. Close the 'done' chan to stop subscription.
func RouteSubscribe(ch chan<- RouteUpdate, done <-chan struct{}) error {
	return routeSubscribeAt(netns.None(), netns.None(), ch, done, nil, false, 0, nil, false)
}

// RouteSubscribeAt works like RouteSubscribe plus it allows the caller
// to choose the network namespace in which to subscribe (ns).
func RouteSubscribeAt(ns netns.NsHandle, ch chan<- RouteUpdate, done <-chan struct{}) error {
	return routeSubscribeAt(ns, netns.None(), ch, done, nil, false, 0, nil, false)
}

// RouteSubscribeOptions contains a set of options to use with
// RouteSubscribeWithOptions.
type RouteSubscribeOptions struct {
	Namespace              *netns.NsHandle
	ErrorCallback          func(error)
	ListExisting           bool
	ReceiveBufferSize      int
	ReceiveBufferForceSize bool
	ReceiveTimeout         *unix.Timeval
}

// RouteSubscribeWithOptions work like RouteSubscribe but enable to
// provide additional options to modify the behavior. Currently, the
// namespace can be provided as well as an error callback.
func RouteSubscribeWithOptions(ch chan<- RouteUpdate, done <-chan struct{}, options RouteSubscribeOptions) error {
	if options.Namespace == nil {
		none := netns.None()
		options.Namespace = &none
	}
	return routeSubscribeAt(*options.Namespace, netns.None(), ch, done, options.ErrorCallback, options.ListExisting,
		options.ReceiveBufferSize, options.ReceiveTimeout, options.ReceiveBufferForceSize)
}

func routeSubscribeAt(newNs, curNs netns.NsHandle, ch chan<- RouteUpdate, done <-chan struct{}, cberr func(error), listExisting bool,
	rcvbuf int, rcvTimeout *unix.Timeval, rcvbufForce bool) error {
	s, err := nl.SubscribeAt(newNs, curNs, unix.NETLINK_ROUTE, unix.RTNLGRP_IPV4_ROUTE, unix.RTNLGRP_IPV6_ROUTE)
	if err != nil {
		return err
	}
	if rcvTimeout != nil {
		if err := s.SetReceiveTimeout(rcvTimeout); err != nil {
			return err
		}
	}
	if rcvbuf != 0 {
		err = s.SetReceiveBufferSize(rcvbuf, rcvbufForce)
		if err != nil {
			return err
		}
	}
	if done != nil {
		go func() {
			<-done
			s.Close()
		}()
	}
	if listExisting {
		req := pkgHandle.newNetlinkRequest(unix.RTM_GETROUTE,
			unix.NLM_F_DUMP)
		infmsg := nl.NewIfInfomsg(unix.AF_UNSPEC)
		req.AddData(infmsg)
		if err := s.Send(req); err != nil {
			return err
		}
	}
	go func() {
		defer close(ch)
		for {
			msgs, from, err := s.Receive()
			if err != nil {
				if cberr != nil {
					cberr(fmt.Errorf("Receive failed: %v",
						err))
				}
				return
			}
			if from.Pid != nl.PidKernel {
				if cberr != nil {
					cberr(fmt.Errorf("Wrong sender portid %d, expected %d", from.Pid, nl.PidKernel))
				}
				continue
			}
			for _, m := range msgs {
				if m.Header.Type == unix.NLMSG_DONE {
					continue
				}
				if m.Header.Type == unix.NLMSG_ERROR {
					error := int32(native.Uint32(m.Data[0:4]))
					if error == 0 {
						continue
					}
					if cberr != nil {
						cberr(fmt.Errorf("error message: %v",
							syscall.Errno(-error)))
					}
					continue
				}
				route, err := deserializeRoute(m.Data)
				if err != nil {
					if cberr != nil {
						cberr(err)
					}
					continue
				}
				ch <- RouteUpdate{
					Type:    m.Header.Type,
					NlFlags: m.Header.Flags & (unix.NLM_F_REPLACE | unix.NLM_F_EXCL | unix.NLM_F_CREATE | unix.NLM_F_APPEND),
					Route:   route,
				}
			}
		}
	}()

	return nil
}

func (p RouteProtocol) String() string {
	switch int(p) {
	case unix.RTPROT_BABEL:
		return "babel"
	case unix.RTPROT_BGP:
		return "bgp"
	case unix.RTPROT_BIRD:
		return "bird"
	case unix.RTPROT_BOOT:
		return "boot"
	case unix.RTPROT_DHCP:
		return "dhcp"
	case unix.RTPROT_DNROUTED:
		return "dnrouted"
	case unix.RTPROT_EIGRP:
		return "eigrp"
	case unix.RTPROT_GATED:
		return "gated"
	case unix.RTPROT_ISIS:
		return "isis"
	// case unix.RTPROT_KEEPALIVED:
	//	return "keepalived"
	case unix.RTPROT_KERNEL:
		return "kernel"
	case unix.RTPROT_MROUTED:
		return "mrouted"
	case unix.RTPROT_MRT:
		return "mrt"
	case unix.RTPROT_NTK:
		return "ntk"
	case unix.RTPROT_OSPF:
		return "ospf"
	case unix.RTPROT_RA:
		return "ra"
	case unix.RTPROT_REDIRECT:
		return "redirect"
	case unix.RTPROT_RIP:
		return "rip"
	case unix.RTPROT_STATIC:
		return "static"
	case unix.RTPROT_UNSPEC:
		return "unspec"
	case unix.RTPROT_XORP:
		return "xorp"
	case unix.RTPROT_ZEBRA:
		return "zebra"
	default:
		return strconv.Itoa(int(p))
	}
}

// genZeroIPNet returns 0.0.0.0/0 or ::/0 for IPv4 or IPv6, otherwise nil
func genZeroIPNet(family int) *net.IPNet {
	var addLen int
	var ip net.IP
	switch family {
	case FAMILY_V4:
		addLen = net.IPv4len
		ip = net.IPv4zero
	case FAMILY_V6:
		addLen = net.IPv6len
		ip = net.IPv6zero
	}
	if addLen != 0 {
		return &net.IPNet{
			IP:   ip,
			Mask: net.CIDRMask(0, 8*addLen),
		}
	}
	return nil
}
