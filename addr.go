package netlink

import (
	"fmt"
	"net"
	"strings"
)

// Addr represents an IP address from netlink. Netlink ip addresses
// include a mask, so it stores the address as a net.IPNet.
type Addr struct {
	*net.IPNet
	Label       string
	Flags       int
	Scope       int
	Peer        *net.IPNet
	Broadcast   net.IP
	PreferedLft int
	ValidLft    int
	LinkIndex   int
}

// String returns $ip/$netmask $label
func (a Addr) String() string {
	return strings.TrimSpace(fmt.Sprintf("%s %s", a.IPNet, a.Label))
}

