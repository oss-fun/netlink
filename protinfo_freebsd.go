package netlink

import (
	"github.com/oss-fun/netlink/nlsyscall"
	"github.com/oss-fun/netlink/nl"
)

func parseProtinfo(infos []nlsyscall.NetlinkRouteAttr) (pi Protinfo) {
	for _, info := range infos {
		switch info.Attr.Type {
		case nl.IFLA_BRPORT_MODE:
			pi.Hairpin = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_GUARD:
			pi.Guard = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_FAST_LEAVE:
			pi.FastLeave = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_PROTECT:
			pi.RootBlock = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_LEARNING:
			pi.Learning = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_UNICAST_FLOOD:
			pi.Flood = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_PROXYARP:
			pi.ProxyArp = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_PROXYARP_WIFI:
			pi.ProxyArpWiFi = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_ISOLATED:
			pi.Isolated = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_NEIGH_SUPPRESS:
			pi.NeighSuppress = byteToBool(info.Value[0])
		}
	}
	return
}
