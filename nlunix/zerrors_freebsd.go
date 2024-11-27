package nlunix

import (
	"syscall"
)

const (
	AF_MPLS                   = 0x1c // Dammy unix.AF_MPLS
	AF_NETLINK                = 38
	ARPHRD_ETHER              = 1
	ARPHRD_IEEE802            = 6
	ARPHRD_ARCNET             = 7
	ARPHRD_DLCI               = 15
	ARPHRD_ATM                = 19
	ARPHRD_IEEE1394           = 24
	ARPHRD_INFINIBAND         = 32
	ARPHRD_RAWHDLC            = 518
	ARPHRD_FRAD               = 770
	ARPHRD_FDDI               = 774
	ARPHRD_SIT                = 776
	ARPHRD_IRDA               = 783
	ARPHRD_FCPP               = 784
	ARPHRD_FCAL               = 785
	ARPHRD_FCPL               = 786
	ARPHRD_FCFABRIC           = 787
	ARPHRD_IEEE802_TR         = 800
	ARPHRD_IEEE80211          = 801
	ARPHRD_IEEE80211_PRISM    = 802
	ARPHRD_IEEE80211_RADIOTAP = 803
	ARPHRD_IEEE802154         = 804
	NETLINK_EXT_ACK           = 0xb
	NETLINK_GET_STRICT_CHK    = 0xc
	NETLINK_ROUTE             = 0x0
	NETLINK_NETFILTER         = 0xc  // not supported
	NETLINK_XFRM              = 0x6  // (not supported) PF_SETKEY
	NLMSG_ERROR               = 0x2
	NLMSG_DONE                = 0x3
	NLMSG_HDRLEN              = 0x10
	NLM_F_ACK                 = 0x4
	NLM_F_ACK_TLVS            = 0x200
	NLM_F_CREATE              = 0x400
	NLM_F_DUMP                = 0x300
	NLM_F_DUMP_INTR           = 0x10
	NLM_F_EXCL                = 0x200
	NLM_F_MULTI               = 0x2
	NLM_F_REQUEST             = 0x1
	RTA_ALIGNTO               = 0x4  // sizeof(uint32_t)
	RTM_GETADDR               = 0x16
	RTM_SETLINK               = 0x13 // not supported
	RTM_NEWADDR               = 0x14 // not supported
	RTM_NEWLINK               = 0x10 
	RTNH_F_PERVASIVE          = 0x2  // not supported
	RTNH_F_ONLINK             = 0x4  // not supported
	RTPROT_BOOT               = 0x3
	SOL_NETLINK               = 0x10e
)

/* Dammy */
const (
	IFF_TUN         = 0x1
	IFF_TAP         = 0x2
	IFF_TUN_EXCL    = 0x8000
	IFF_ONE_QUEUE   = 0x2000 
	IFF_VNET_HDR    = 0x4000
	IFF_NO_PI       = 0x1000
	IFF_MULTI_QUEUE = 0x100

	TUNSETIFF     = 0x400454ca
	TUNSETPERSIST = 0x400454ca
)

/* zerrors_freebsd_arch.go */
const (
	SO_RCVBUF      = 0x1002
	SO_RCVBUFFORCE = 0x21 // not supported
)

const (
	E2BIG       = syscall.Errno(0x7)
	EACCES      = syscall.Errno(0xd)
	EAGAIN      = syscall.Errno(0xb)
	EBADF       = syscall.Errno(0x9)
	EBUSY       = syscall.Errno(0x10)
	ECHILD      = syscall.Errno(0xa)
	EDOM        = syscall.Errno(0x21)
	EEXIST      = syscall.Errno(0x11)
	EFAULT      = syscall.Errno(0xe)
	EFBIG       = syscall.Errno(0x1b)
	EINTR       = syscall.Errno(0x4)
	EINVAL      = syscall.Errno(0x16)
	EIO         = syscall.Errno(0x5)
	EISDIR      = syscall.Errno(0x15)
	EMFILE      = syscall.Errno(0x18)
	EMLINK      = syscall.Errno(0x1f)
	ENFILE      = syscall.Errno(0x17)
	ENODEV      = syscall.Errno(0x13)
	ENOENT      = syscall.Errno(0x2)
	ENOEXEC     = syscall.Errno(0x8)
	ENOMEM      = syscall.Errno(0xc)
	ENOSPC      = syscall.Errno(0x1c)
	ENOTBLK     = syscall.Errno(0xf)
	ENOTDIR     = syscall.Errno(0x14)
	ENOTTY      = syscall.Errno(0x19)
	ENXIO       = syscall.Errno(0x6)
	EPERM       = syscall.Errno(0x1)
	EPIPE       = syscall.Errno(0x20)
	ERANGE      = syscall.Errno(0x22)
	EROFS       = syscall.Errno(0x1e)
	ESPIPE      = syscall.Errno(0x1d)
	ESRCH       = syscall.Errno(0x3)
	ETXTBSY     = syscall.Errno(0x1a)
	EWOULDBLOCK = syscall.Errno(0xb)
	EXDEV       = syscall.Errno(0x12)
)

