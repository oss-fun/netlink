package main

import (
	"fmt"
	"net"
	"log"
	"os"

	"github.com/oss-fun/netlink"
)

func main() {
	fmt.Println("FreeBSD netlink-addradd")
	if len(os.Args) < 3 {
		fmt.Println("Usage: netlink-addradd <interface-name> <ipaddr>")
		os.Exit(1)
	}

	ifaceName := os.Args[1]
	ipAddr := os.Args[2]

	link, err := netlink.LinkByName(ifaceName)
	if err != nil {
		log.Fatalf("Failed to get interface %s: %v", ifaceName, err)
	}

	ipNet := &net.IPNet{
		IP:   net.ParseIP(ipAddr),
		Mask: net.CIDRMask(24, 32),
	}
	addr := &netlink.Addr{IPNet: ipNet}

	if err := netlink.AddrAdd(link, addr); err != nil {
		log.Fatalf("Failed to add address %s to interface %s: %v", ipAddr, ifaceName, err)
	}

	fmt.Printf("Successfully added %s to %s\n", ipAddr, ifaceName)
}
