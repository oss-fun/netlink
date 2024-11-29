package nlunix

import (
	"syscall"
)

const (
	AF_BRIDGE                 = 0x7    // Dammy
	AF_MPLS                   = 0x1c   // Dammy unix.AF_MPLS
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
	NETLINK_NETFILTER         = 0xc    // not supported
	NETLINK_XFRM              = 0x6    // (not supported) PF_SETKEY
	NLMSG_ERROR               = 0x2
	NLMSG_DONE                = 0x3
	NLMSG_HDRLEN              = 0x10
	NLA_F_NESTED              = 0x8000
	NLM_F_ACK                 = 0x4
	NLM_F_ACK_TLVS            = 0x200
	NLM_F_APPEND              = 0x800
	NLM_F_CREATE              = 0x400
	NLM_F_DUMP                = 0x300
	NLM_F_DUMP_INTR           = 0x10
	NLM_F_EXCL                = 0x200
	NLM_F_MULTI               = 0x2
	NLM_F_REPLACE             = 0x100
	NLM_F_REQUEST             = 0x1
	RTAX_MTU                  = 0x2
	RTAX_WINDOW               = 0x3    // not supported
	RTAX_RTT                  = 0x4    // not supported
	RTAX_RTTVAR               = 0x5    // not supported
	RTAX_SSTHRESH             = 0x6    // not supported
	RTAX_CWND                 = 0x7    // not supported
	RTAX_ADVMSS               = 0x8    // not supported
	RTAX_REORDERING           = 0x9    // not supported
	RTAX_HOPLIMIT             = 0xa    // not supported
	RTAX_INITCWND             = 0xb    // not supported
	RTAX_FEATURES             = 0xc    // not supported
	RTAX_RTO_MIN              = 0xd    // not supported
	RTAX_INITRWND             = 0xe    // not supported
	RTAX_QUICKACK             = 0xf    // not supported
	RTAX_CC_ALGO              = 0x10   // not supported
	RTAX_FASTOPEN_NO_COOKIE   = 0x11   // not supported
	RTA_ALIGNTO               = 0x4    // sizeof(uint32_t)
	RTM_DELLINK               = 0x11
	RTM_DELROUTE              = 0x19
	RTM_DELRULE               = 0x21   // not supported
	RTM_DELNEIGH              = 0x1d
	RTM_GETADDR               = 0x16
	RTM_GETLINK               = 0x12
	RTM_GETNEIGH              = 0x1e
	RTM_GETROUTE              = 0x1a
	RTM_GETRULE               = 0x22   // not supported
	RTM_SETLINK               = 0x13   // not supported
	RTM_NEWADDR               = 0x14   // not supported
	RTM_NEWLINK               = 0x10
	RTM_NEWRULE               = 0x20   // not supported
	RTM_NEWNEIGH              = 0x1c
	RTM_NEWROUTE              = 0x18
	RTM_F_CLONED              = 0x200  // not supported
	RTM_F_FIB_MATCH           = 0x2000 // not supported
	RTM_F_LOOKUP_TABLE        = 0x1000 // not supported
	RTNH_F_PERVASIVE          = 0x2    // not supported
	RTNH_F_ONLINK             = 0x4    // not supported
	RTPROT_BABEL              = 0x2a
	RTPROT_BGP                = 0xba
	RTPROT_BIRD               = 0xc
	RTPROT_BOOT               = 0x3
	RTPROT_DHCP               = 0x10
	RTPROT_DNROUTED           = 0xd
	RTPROT_EIGRP              = 0xc0
	RTPROT_GATED              = 0x8
	RTPROT_ISIS               = 0xbb
	RTPROT_KEEPALIVED         = 0x12
	RTPROT_KERNEL             = 0x2
	RTPROT_MROUTED            = 0x11
	RTPROT_MRT                = 0xa
	RTPROT_NTK                = 0xf
	RTPROT_OPENR              = 0x63
	RTPROT_OSPF               = 0xbc
	RTPROT_RA                 = 0x9
	RTPROT_REDIRECT           = 0x1
	RTPROT_RIP                = 0xbd
	RTPROT_STATIC             = 0x4
	RTPROT_UNSPEC             = 0x0
	RTPROT_XORP               = 0xe
	RTPROT_ZEBRA              = 0xb	
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
	IFF_PERSIST     = 0x800

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

