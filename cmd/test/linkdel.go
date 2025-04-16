package main

import (
	"fmt"
	"log"
	"os"

	"github.com/oss-fun/netlink"
)

func main() {
	fmt.Println("FreeBSD netlink-linkdel")
	if len(os.Args) < 2 {
		fmt.Println("Usage: netlink-linkdel <interface-name>")
		os.Exit(1)
	}

	ifaceName := os.Args[1]

	link, err := netlink.LinkByName(ifaceName)
	if err != nil {
		log.Fatalf("Failed to find interface: %v", err)
	}

	if err := netlink.LinkDel(link); err != nil {
		log.Fatalf("Failed to delete interface: %v", err)
	}

	fmt.Printf("Interface %s deleted successfully\n", ifaceName)
}
