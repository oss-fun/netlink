package main

import (
	"os"
	"log"
	"fmt"
	"strconv"

	"github.com/oss-fun/netlink"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: netlink-setnsfd <ifName> <jid>\n")
		os.Exit(1)
	}

	ifName := os.Args[1]
	jid := os.Args[2]

	link, err := netlink.LinkByName(ifName)
	if err != nil {
		log.Fatalf("インターフェース %s を取得できません: %v", ifName, err)
	}

	jid_int, err := strconv.Atoi(jid)
	if err != nil {
		log.Fatalf("strconv error: %v", err)
	}
	if err = netlink.LinkSetNsFd(link, jid_int); err != nil {
		log.Fatalf("インターフェース（%v）をjail（jid: %v）に追加できません: %v", ifName, jid, err)
	}
}
