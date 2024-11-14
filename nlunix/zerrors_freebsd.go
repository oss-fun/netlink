package nlunix

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
	NETLINK_ROUTE             = 0x0
	NETLINK_NETFILTER         = 0xc  // not supported
	NETLINK_XFRM              = 0x6  // (not supported) PF_SETKEY
	NLMSG_ERROR               = 0x2
	NLMSG_DONE                = 0x3
	NLMSG_HDRLEN              = 0x10
	NLM_F_ACK_TLVS            = 0x200
	NLM_F_DUMP_INTR           = 0x10
	NLM_F_MULTI               = 0x2
	NLM_F_REQUEST             = 0x1
	RTA_ALIGNTO               = 0x4  // sizeof(uint32_t)
	RTPROT_BOOT               = 0x3
	SOL_NETLINK               = 0x10e
)

/* zerrors_freebsd_arch.go */
const (
	SO_RCVBUF      = 0x1002
	SO_RCVBUFFORCE = 0x21 // not supported
)
