package nlunix

type RawSockaddrNetlink struct {
	Family uint16
	Pad    uint16
	Pid    uint32
	Groups uint32
}

const (
	RT_SCOPE_NOWHERE   = 0xff
	RT_SCOPE_UNIVERSE  = 0x0
	RT_TABLE_MAIN      = 0x0 // RT_DEFAULT_FIB
	RT_TABLE_UNSPEC    = 0xFFFFFFFF // RT_ALL_FIBS
	SizeofNlMsghdr     = 0x10
	SizeofRtAttr       = 0x4
	SizeofIfInfomsg    = 0x10
	SizeofIfAddrmsg    = 0x8
	SizeofIfaCacheinfo = 0x10
)

type NlMsghdr struct {
	Len   uint32
	Type  uint16
	Flags uint16
	Seq   uint32
	Pid   uint32
}

type RtAttr struct {
	Len  uint16
	Type uint16
}

type IfInfomsg struct {
	Family uint8
	_      uint8
	Type   uint16
	Index  int32
	Flags  uint32
	Change uint32
}

type IfAddrmsg struct {
	Family    uint8
	Prefixlen uint8
	Flags     uint8
	Scope     uint8
	Index     uint32
}

type IfaCacheinfo struct {
	Prefered uint32
	Valid    uint32
	Cstamp   uint32
	Tstamp   uint32
}

type RtMsg struct {
	Family   uint8
	Dst_len  uint8
	Src_len  uint8
	Tos      uint8
	Table    uint8
	Protocol uint8
	Scope    uint8
	Type     uint8
	Flags    uint32
}

type RtNexthop struct {
	Len     uint16
	Flags   uint8
	Hops    uint8
	Ifindex int32
}
