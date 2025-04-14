package main

import (
	"fmt"
	"log"
	"os"

	"github.com/oss-fun/netlink"
)

func main() {
	fmt.Println("FreeBSD netlink-linkadd")
	if len(os.Args) < 3 {
		fmt.Println("Usage: netlink-linkadd <interface-name> <peer-interface-name>")
		os.Exit(1)
	}

	ifaceName := os.Args[1]
	peerIfaceName := os.Args[2]

	// ダミーインターフェースを作成
	link := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name: ifaceName,
		},
		PeerName: peerIfaceName,
	}

	if err := netlink.LinkAdd(link); err != nil {
		log.Fatalf("Failed to add veth pair: %v", err)
	}

	fmt.Printf("Veth pair %s <-> %s created successfully.\n", ifaceName, peerIfaceName)
}
