package netlink

type ethtoolStats struct {
	cmd    uint32
	nStats uint32
	// Followed by nStats * []uint64.
}

