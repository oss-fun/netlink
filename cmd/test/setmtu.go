package main

import (
    "fmt"
    "log"
    "os"
    "strconv"

    "github.com/oss-fun/netlink"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: set-mtu <interface-name> <mtu>")
        os.Exit(1)
    }

    ifaceName := os.Args[1]
    mtuStr := os.Args[2]

    mtu, err := strconv.Atoi(mtuStr)
    if err != nil {
        log.Fatalf("MTU値が正しくありません: %v", err)
    }

    // インターフェースを取得
    link, err := netlink.LinkByName(ifaceName)
    if err != nil {
        log.Fatalf("インターフェース %s が見つかりません: %v", ifaceName, err)
    }

    // MTU を設定
    if err := netlink.LinkSetMTU(link, mtu); err != nil {
        log.Fatalf("MTU の設定に失敗しました: %v", err)
    }

    fmt.Printf("インターフェース %s の MTU を %d に設定しました\n", ifaceName, mtu)
}

