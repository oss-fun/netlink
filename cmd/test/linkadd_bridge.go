package main

import (
	"fmt"
	"log"

	"github.com/oss-fun/netlink"
)

func main() {
	// 作成するブリッジインターフェースの名前
	bridgeName := "br0"

	// ブリッジデバイスを作成
	link := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: bridgeName,
		},
	}

	// インターフェースを追加
	if err := netlink.LinkAdd(link); err != nil {
		log.Fatalf("Failed to add bridge: %v", err)
	}

	fmt.Printf("Bridge %s created successfully.\n", bridgeName)
}
