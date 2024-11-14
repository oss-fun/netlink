package nl

import (
	"unsafe"

	"github.com/oss-fun/netlink/nlunix"
)

type IfAddrmsg struct {
	nlunix.IfAddrmsg
}

func NewIfAddrmsg(family int) *IfAddrmsg {
	return &IfAddrmsg{
		IfAddrmsg: nlunix.IfAddrmsg{
			Family: uint8(family),
		},
	}
}

func DeserializeIfAddrmsg(b []byte) *IfAddrmsg {
	return (*IfAddrmsg)(unsafe.Pointer(&b[0:nlunix.SizeofIfAddrmsg][0]))
}

func (msg *IfAddrmsg) Serialize() []byte {
	return (*(*[nlunix.SizeofIfAddrmsg]byte)(unsafe.Pointer(msg)))[:]
}

func (msg *IfAddrmsg) Len() int {
	return nlunix.SizeofIfAddrmsg
}

type IfaCacheInfo struct {
	nlunix.IfaCacheinfo
}

func (msg *IfaCacheInfo) Len() int {
	return nlunix.SizeofIfaCacheinfo
}

func DeserializeIfaCacheInfo(b []byte) *IfaCacheInfo {
	return (*IfaCacheInfo)(unsafe.Pointer(&b[0:nlunix.SizeofIfaCacheinfo][0]))
}

func (msg *IfaCacheInfo) Serialize() []byte {
	return (*(*[nlunix.SizeofIfaCacheinfo]byte)(unsafe.Pointer(msg)))[:]
}
