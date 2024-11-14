package nl

import (
	"unsafe"

	"github.com/oss-fun/netlink/nlunix"
)

type RtMsg struct {
	nlunix.RtMsg
}

func NewRtMsg() *RtMsg {
	return &RtMsg{
		RtMsg: nlunix.RtMsg{
			Table:    nlunix.RT_TABLE_MAIN,
			Scope:    nlunix.RT_SCOPE_UNIVERSE,
			Protocol: nlunix.RTPROT_BOOT,
			Type:     nlunix.RTN_UNICAST,
		},
	}
}

func NewRtDelMsg() *RtMsg {
	return &RtMsg{
		RtMsg: nlunix.RtMsg{
			Table: nlunix.RT_TABLE_MAIN,
			Scope: nlunix.RT_SCOPE_NOWHERE,
		},
	}
}

func (msg *RtMsg) Len() int {
	return nlunix.SizeofRtMsg
}

func DeserializeRtMsg(b []byte) *RtMsg {
	return (*RtMsg)(unsafe.Pointer(&b[0:nlunix.SizeofRtMsg][0]))
}

func (msg *RtMsg) Serialize() []byte {
	return (*(*[nlunix.SizeofRtMsg]byte)(unsafe.Pointer(msg)))[:]
}

type RtNexthop struct {
	nlunix.RtNexthop
	Children []NetlinkRequestData
}

func DeserializeRtNexthop(b []byte) *RtNexthop {
	return &RtNexthop{
		RtNexthop: *((*nlunix.RtNexthop)(unsafe.Pointer(&b[0:nlunix.SizeofRtNexthop][0]))),
	}
}

func (msg *RtNexthop) Len() int {
	if len(msg.Children) == 0 {
		return nlunix.SizeofRtNexthop
	}

	l := 0
	for _, child := range msg.Children {
		l += rtaAlignOf(child.Len())
	}
	l += nlunix.SizeofRtNexthop
	return rtaAlignOf(l)
}

func (msg *RtNexthop) Serialize() []byte {
	length := msg.Len()
	msg.RtNexthop.Len = uint16(length)
	buf := make([]byte, length)
	copy(buf, (*(*[nlunix.SizeofRtNexthop]byte)(unsafe.Pointer(msg)))[:])
	next := rtaAlignOf(nlunix.SizeofRtNexthop)
	if len(msg.Children) > 0 {
		for _, child := range msg.Children {
			childBuf := child.Serialize()
			copy(buf[next:], childBuf)
			next += rtaAlignOf(len(childBuf))
		}
	}
	return buf
}

