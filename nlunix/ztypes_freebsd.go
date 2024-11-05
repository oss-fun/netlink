package nlunix

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
