package main

import (
    "fmt"
    "log"

    "github.com/oss-fun/netlink"
)

func main() {
    // すべてのルート（IPv4 + IPv6）を取得
    routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
    if err != nil {
        log.Fatalf("ルーティングテーブルの取得に失敗しました: %v", err)
    }

    fmt.Println("ルーティングテーブル一覧:")
    for _, route := range routes {
        dst := "default"
        if route.Dst != nil {
            dst = route.Dst.String()
        }

        gw := "なし"
        if route.Gw != nil {
            gw = route.Gw.String()
        }

        iface := "不明"
        if link, err := netlink.LinkByIndex(route.LinkIndex); err == nil {
            iface = link.Attrs().Name
        }

        fmt.Printf(" - 宛先: %-20s ゲートウェイ: %-15s インターフェース: %s\n", dst, gw, iface)
    }
}

