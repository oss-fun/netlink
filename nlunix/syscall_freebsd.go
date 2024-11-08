package nlunix

type SockaddrNetlink struct {
	Family uint16
	Pad    uint16
	Pid    uint16
	Groups uint32
	raw    RawSockaddrNetlink
}
