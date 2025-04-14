package netlink

import (
	"fmt"
)

func atob(a []byte) ([]byte, error) {
	if len(a) < 2 {
		return nil, fmt.Errorf("invalid interface name: %q", a)
	}
	if a[len(a)-1] != 'a' {
		return nil, fmt.Errorf("interface name does not end with 'a': %q", a)
	}

	b := make([]byte, len(a))
	copy(b, a)
	b[len(b)-1] = 'b'
	return b, nil
}
