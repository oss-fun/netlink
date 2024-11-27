package nlunix

type RawSockaddrNetlink struct {
	Family uint16
	Pad    uint16
	Pid    uint32
	Groups uint32
}

type RawSockaddrInet4 struct {
	Family uint16
	Port   uint16
	Addr   [4]byte /* in_addr */
	Zero   [8]uint8
}

type RawSockaddrInet6 struct {
	Family   uint16
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
}

type RawSockaddrUnix struct {
	Len    uint8
	Family uint8
	Path   [104]int8
}

type RawSockaddrDatalink struct {
	Len    uint8
	Family uint8
	Index  uint16
	Type   uint8
	Nlen   uint8
	Alen   uint8
	Slen   uint8
	Data   [46]int8
}

type RawSockaddr struct {
	Family uint16
	Data   [14]int8
}

type RawSockaddrAny struct {
	Addr RawSockaddr
	Pad  [96]int8
}

type _Socklen uint32

const (
	SizeofSockaddrInet4    = 0x10
	SizeofSockaddrInet6    = 0x1c
	SizeofSockaddrAny      = 0x70
	SizeofSockaddrUnix     = 0x6a
	SizeofSockaddrDatalink = 0x36
	SizeofSockaddrNetlink  = 0xc
)

const (
	IFA_UNSPEC         = 0x0
	IFA_ADDRESS        = 0x1
	IFA_LOCAL          = 0x2
	IFA_LABEL          = 0x3
	IFA_BROADCAST      = 0x4
	IFA_ANYCAST        = 0x5
	IFA_CACHEINFO      = 0x6
	IFA_MULTICAST      = 0x7
	IFA_FLAGS          = 0x8
	RT_SCOPE_UNIVERSE  = 0x0
	RT_SCOPE_SITE      = 0xc8
	RT_SCOPE_LINK      = 0xfd
	RT_SCOPE_HOST      = 0xfe
	RT_SCOPE_NOWHERE   = 0xff
	RT_TABLE_MAIN      = 0x0 // RT_DEFAULT_FIB
	RT_TABLE_UNSPEC    = 0xFFFFFFFF // RT_ALL_FIBS
	RTN_UNICAST        = 0x1
	SizeofNlMsghdr     = 0x10
	SizeofRtAttr       = 0x4
	SizeofIfInfomsg    = 0x10
	SizeofIfAddrmsg    = 0x8
	SizeofIfaCacheinfo = 0x10
	SizeofRtMsg        = 0xc
	SizeofRtNexthop    = 0x8
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

const (
	IFLA_ADDRESS           = 0x1
	IFLA_IFNAME            = 0x3
	IFLA_MTU               = 0x4
	IFLA_LINK              = 0x5
	IFLA_TXQLEN            = 0xd
	IFLA_LINKINFO          = 0x12
	IFLA_NET_NS_PID        = 0x13
	IFLA_IFALIAS           = 0x14
	IFLA_GROUP             = 0x1b
	IFLA_NET_NS_FD         = 0x1c
	IFLA_NUM_TX_QUEUES     = 0x1f
	IFLA_NUM_RX_QUEUES     = 0x20
	IFLA_GSO_MAX_SEGS      = 0x28
	IFLA_GSO_MAX_SIZE      = 0x29
	IFLA_GRO_MAX_SIZE      = 0x3a
	IFLA_GSO_IPV4_MAX_SIZE = 0x3f
	IFLA_GRO_IPV4_MAX_SIZE = 0x40
)
