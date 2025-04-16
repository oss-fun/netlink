package main

import (
	"fmt"
	"log"
	"os"

	"github.com/oss-fun/netlink"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: netlink-linksetup <interface_name>")
		os.Exit(1)
	}
	
	ifName := os.Args[1]

	// インターフェースの取得
	link, err := netlink.LinkByName(ifName)
	if err != nil {
		log.Fatalf("インターフェース %s を取得できません: %v", ifName, err)
	}

	// インターフェースを UP に設定
	if err := netlink.LinkSetUp(link); err != nil {
		log.Fatalf("インターフェース %s を UP にできません: %v", ifName, err)
	}

	fmt.Printf("インターフェース %s を UP にしました\n", ifName)
}
