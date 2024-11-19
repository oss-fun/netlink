package nl

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"testing"

	"github.com/oss-fun/netlink/nlunix"
)

func (msg *RtMsg) write(b []byte) {
	native := NativeEndian()
	b[0] = msg.Family
	b[1] = msg.Dst_len
	b[2] = msg.Src_len
	b[3] = msg.Tos
	b[4] = msg.Table
	b[5] = msg.Protocol
	b[6] = msg.Scope
	b[7] = msg.Type
	native.PutUint32(b[8:12], msg.Flags)
}

func (msg *RtMsg) serializeSafe() []byte {
	len := nlunix.SizeofRtMsg
	b := make([]byte, len)
	msg.write(b)
	return b
}

func deserializeRtMsgSafe(b []byte) *RtMsg {
	var msg = RtMsg{}
	binary.Read(bytes.NewReader(b[0:nlunix.SizeofRtMsg]), NativeEndian(), &msg)
	return &msg
}

func TestRtMsgDeserializeSerialize(t *testing.T) {
	var orig = make([]byte, nlunix.SizeofRtMsg)
	rand.Read(orig)
	safemsg := deserializeRtMsgSafe(orig)
	msg := DeserializeRtMsg(orig)
	testDeserializeSerialize(t, orig, safemsg, msg)
}

func TestDeserializeRtNexthop(t *testing.T) {
	buf := make([]byte, nlunix.SizeofRtNexthop+64)
	native := NativeEndian()
	native.PutUint16(buf[0:2], nlunix.SizeofRtNexthop)
	buf[2] = 17
	buf[3] = 1
	native.PutUint32(buf[4:8], 1234)

	msg := DeserializeRtNexthop(buf)
	safemsg := &RtNexthop{
		nlunix.RtNexthop{
			Len:     nlunix.SizeofRtNexthop,
			Flags:   17,
			Hops:    1,
			Ifindex: 1234,
		},
		nil,
	}
	if msg.Len() != safemsg.Len() || msg.Flags != safemsg.Flags || msg.Hops != safemsg.Hops || msg.Ifindex != safemsg.Ifindex {
		t.Fatal("Deserialization failed.\nIn:", buf, "\nOut:", msg, "\n", msg.Serialize(), "\nExpected:", safemsg, "\n", safemsg.Serialize())
	}
}
