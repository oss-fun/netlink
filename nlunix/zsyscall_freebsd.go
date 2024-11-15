package nlunix

import (
	"unsafe"	
	
	"golang.org/x/sys/unix"
)

func bind(s int, addr unsafe.Pointer, addrlen _Socklen) (err error) {
	_, _, e1 := unix.Syscall(unix.SYS_BIND, uintptr(s), uintptr(addr), uintptr(addrlen))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}
