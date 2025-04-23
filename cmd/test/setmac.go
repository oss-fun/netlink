package main

import (
    "fmt"
    "log"
    "net"
    "os"

    "github.com/oss-fun/netlink"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: set-mac <interface-name> <mac-address>")
        os.Exit(1)
    }

    ifaceName := os.Args[1]
    macAddrStr := os.Args[2]

    // MAC アドレスをパース
    hwAddr, err := net.ParseMAC(macAddrStr)
    if err != nil {
        log.Fatalf("MAC アドレスの形式が正しくありません: %v", err)
    }

    // インターフェースを取得
    link, err := netlink.LinkByName(ifaceName)
    if err != nil {
        log.Fatalf("インターフェース %s の取得に失敗: %v", ifaceName, err)
    }

    // インターフェースを一旦 DOWN にする必要がある
    if err := netlink.LinkSetDown(link); err != nil {
        log.Fatalf("インターフェースを DOWN にできませんでした: %v", err)
    }

    // MAC アドレスを設定
    if err := netlink.LinkSetHardwareAddr(link, hwAddr); err != nil {
        log.Fatalf("MAC アドレスの設定に失敗しました: %v", err)
    }

    // インターフェースを再度 UP にする
    if err := netlink.LinkSetUp(link); err != nil {
        log.Fatalf("インターフェースを UP にできませんでした: %v", err)
    }

    fmt.Printf("インターフェース %s の MAC アドレスを %s に変更しました\n", ifaceName, hwAddr)
}

