package nl

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"testing"

	"github.com/oss-fun/netlink/nlunix"
)

func (msg *IfAddrmsg) write(b []byte) {
	native := NativeEndian()
	b[0] = msg.Family
	b[1] = msg.Prefixlen
	b[2] = msg.Flags
	b[3] = msg.Scope
	native.PutUint32(b[4:8], msg.Index)
}

func (msg *IfAddrmsg) serializeSafe() []byte {
	len := nlunix.SizeofIfAddrmsg
	b := make([]byte, len)
	msg.write(b)
	return b
}

func deserializeIfAddrmsgSafe(b []byte) *IfAddrmsg {
	var msg = IfAddrmsg{}
	binary.Read(bytes.NewReader(b[0:nlunix.SizeofIfAddrmsg]), NativeEndian(), &msg)
	return &msg
}

func TestIfAddrmsgDeserializeSerialize(t *testing.T) {
	var orig = make([]byte, nlunix.SizeofIfAddrmsg)
	rand.Read(orig)
	safemsg := deserializeIfAddrmsgSafe(orig)
	msg := DeserializeIfAddrmsg(orig)
	testDeserializeSerialize(t, orig, safemsg, msg)
}

func (msg *IfaCacheInfo) write(b []byte) {
	native := NativeEndian()
	native.PutUint32(b[0:4], uint32(msg.Prefered))
	native.PutUint32(b[4:8], uint32(msg.Valid))
	native.PutUint32(b[8:12], uint32(msg.Cstamp))
	native.PutUint32(b[12:16], uint32(msg.Tstamp))
}

func (msg *IfaCacheInfo) serializeSafe() []byte {
	length := nlunix.SizeofIfaCacheinfo
	b := make([]byte, length)
	msg.write(b)
	return b
}

func deserializeIfaCacheInfoSafe(b []byte) *IfaCacheInfo {
	var msg = IfaCacheInfo{}
	binary.Read(bytes.NewReader(b[0:nlunix.SizeofIfaCacheinfo]), NativeEndian(), &msg)
	return &msg
}

func TestIfaCacheInfoDeserializeSerialize(t *testing.T) {
	var orig = make([]byte, nlunix.SizeofIfaCacheinfo)
	rand.Read(orig)
	safemsg := deserializeIfaCacheInfoSafe(orig)
	msg := DeserializeIfaCacheInfo(orig)
	testDeserializeSerialize(t, orig, safemsg, msg)
}
