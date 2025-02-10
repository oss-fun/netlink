// +build freebsd

package netlink

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/oss-fun/netlink/nl"
	"github.com/oss-fun/vnet"
	"golang.org/x/sys/unix"
)

func TestHandleCreateClose(t *testing.T) {
	h, err := NewHandle()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range nl.SupportedNlFamilies {
		sh, ok := h.sockets[f]
		if !ok {
			t.Fatalf("Handle socket(s) for family %d was not created", f)
		}
		if sh.Socket == nil {
			t.Fatalf("Socket for family %d was not created", f)
		}
	}

	h.Close()
	if h.sockets != nil {
		t.Fatalf("Handle socket(s) were not closed")
	}
}

func TestHandleCreateNetns(t *testing.T) {
	skipUnlessRoot(t)

	id := make([]byte, 4)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		t.Fatal(err)
	}
	ifName := "dummy-" + hex.EncodeToString(id)

	// Create an handle on the current netns
	curNs, err := vnet.Get()
	if err != nil {
		t.Fatal(err)
	}
	defer curNs.Close()

	ch, err := NewHandleAt(curNs)
	if err != nil {
		t.Fatal(err)
	}
	defer ch.Close()

	// Create an handle on a custom netns
	newNs, err := vnet.New()
	if err != nil {
		t.Fatal(err)
	}
	defer newNs.Close()

	nh, err := NewHandleAt(newNs)
	if err != nil {
		t.Fatal(err)
	}
	defer nh.Close()

	// Create an interface using the current handle
	err = ch.LinkAdd(&Dummy{LinkAttrs{Name: ifName}})
	if err != nil {
		t.Fatal(err)
	}
	l, err := ch.LinkByName(ifName)
	if err != nil {
		t.Fatal(err)
	}
	if l.Type() != "dummy" {
		t.Fatalf("Unexpected link type: %s", l.Type())
	}

	// Verify the new handle cannot find the interface
	ll, err := nh.LinkByName(ifName)
	if err == nil {
		t.Fatalf("Unexpected link found on netns %s: %v", newNs, ll)
	}

	// Move the interface to the new netns
	err = ch.LinkSetNsFd(l, int(newNs))
	if err != nil {
		t.Fatal(err)
	}

	// Verify new netns handle can find the interface while current cannot
	ll, err = nh.LinkByName(ifName)
	if err != nil {
		t.Fatal(err)
	}
	if ll.Type() != "dummy" {
		t.Fatalf("Unexpected link type: %s", ll.Type())
	}
	ll, err = ch.LinkByName(ifName)
	if err == nil {
		t.Fatalf("Unexpected link found on netns %s: %v", curNs, ll)
	}
}

func TestHandleTimeout(t *testing.T) {
	h, err := NewHandle()
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()

	for _, sh := range h.sockets {
		verifySockTimeVal(t, sh.Socket.GetFd(), unix.Timeval{Sec: 0, Usec: 0})
	}

	h.SetSocketTimeout(2*time.Second + 8*time.Millisecond)

	for _, sh := range h.sockets {
		verifySockTimeVal(t, sh.Socket.GetFd(), unix.Timeval{Sec: 2, Usec: 8000})
	}
}

func TestHandleReceiveBuffer(t *testing.T) {
	h, err := NewHandle()
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()
	if err := h.SetSocketReceiveBufferSize(65536, false); err != nil {
		t.Fatal(err)
	}
	sizes, err := h.GetSocketReceiveBufferSize()
	if err != nil {
		t.Fatal(err)
	}
	if len(sizes) != len(h.sockets) {
		t.Fatalf("Unexpected number of socket buffer sizes: %d (expected %d)",
			len(sizes), len(h.sockets))
	}
	for _, s := range sizes {
		if s < 65536 || s > 2*65536 {
			t.Fatalf("Unexpected socket receive buffer size: %d (expected around %d)",
				s, 65536)
		}
	}
}

func verifySockTimeVal(t *testing.T, fd int, tv unix.Timeval) {
	var (
		tr unix.Timeval
		v  = uint32(0x10)
	)
	_, _, errno := unix.Syscall6(unix.SYS_GETSOCKOPT, uintptr(fd), unix.SOL_SOCKET, unix.SO_SNDTIMEO, uintptr(unsafe.Pointer(&tr)), uintptr(unsafe.Pointer(&v)), 0)
	if errno != 0 {
		t.Fatal(errno)
	}

	if tr.Sec != tv.Sec || tr.Usec != tv.Usec {
		t.Fatalf("Unexpected timeout value read: %v. Expected: %v", tr, tv)
	}

	_, _, errno = unix.Syscall6(unix.SYS_GETSOCKOPT, uintptr(fd), unix.SOL_SOCKET, unix.SO_RCVTIMEO, uintptr(unsafe.Pointer(&tr)), uintptr(unsafe.Pointer(&v)), 0)
	if errno != 0 {
		t.Fatal(errno)
	}

	if tr.Sec != tv.Sec || tr.Usec != tv.Usec {
		t.Fatalf("Unexpected timeout value read: %v. Expected: %v", tr, tv)
	}
}

var (
	iter      = 10
	numThread = uint32(4)
	prefix    = "iface"
	handle1   *Handle
	handle2   *Handle
	ns1       vnet.VjHandle
	ns2       vnet.VjHandle
	done      uint32
	initError error
	once      sync.Once
)

func initParallel() {
	ns1, initError = vnet.New()
	if initError != nil {
		return
	}
	handle1, initError = NewHandleAt(ns1)
	if initError != nil {
		return
	}
	ns2, initError = vnet.New()
	if initError != nil {
		return
	}
	handle2, initError = NewHandleAt(ns2)
	if initError != nil {
		return
	}
}

func parallelDone() {
	atomic.AddUint32(&done, 1)
	if done == numThread {
		if ns1.IsOpen() {
			ns1.Close()
		}
		if ns2.IsOpen() {
			ns2.Close()
		}
		if handle1 != nil {
			handle1.Close()
		}
		if handle2 != nil {
			handle2.Close()
		}
	}
}

