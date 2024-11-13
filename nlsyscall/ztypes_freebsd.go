package nlsyscall

const (
	SizeofNlMsghdr = 0x10
	SizeofRtAttr   = 0x4
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
 
