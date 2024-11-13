package nlsyscall

import (
	"unsafe"
)

func nlmAlignOf(msglen int) int {
	return (msglen + NLMSG_ALIGNTO - 1) & ^(NLMSG_ALIGNTO - 1)
}

type NetlinkMessage struct {
	Header NlMsghdr
	Data   []byte
}

func ParseNetlinkMessage(b []byte) ([]NetlinkMessage, error) {
	var msgs []NetlinkMessage
	for len(b) >= NLMSG_HDRLEN {
		h, dbuf, dlen, err := netlinkMessageHeaderAndData(b)
		if err != nil {
			return nil, err
		}
		m := NetlinkMessage{Header: *h, Data: dbuf[:int(h.Len)-NLMSG_HDRLEN]}
		msgs = append(msgs, m)
		b = b[dlen:]
	}
	return msgs, nil
}

func netlinkMessageHeaderAndData(b []byte) (*NlMsghdr, []byte, int, error) {
	h := (*NlMsghdr)(unsafe.Pointer(&b[0]))
	l := nlmAlignOf(int(h.Len))
	if int(h.Len) < NLMSG_HDRLEN || l > len(b) {
		return nil, nil, 0, EINVAL
	}
	return h, b[NLMSG_HDRLEN:], l, nil
}

type NetlinkRouteAttr struct {
	Attr  RtAttr
	Value []byte
}
