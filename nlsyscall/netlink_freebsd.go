package nlsyscall

import (
	"github.com/oss-fun/netlink/nlunix"
)

type NetlinkMessage struct {
	Header nlunix.NlMsghdr
	Data   []byte
}

type NetlinkRouteAttr struct {
	Attr  nlunix.RtAttr
	Value []byte
}
