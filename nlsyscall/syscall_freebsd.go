package nlsyscall

import (
	errorspkg "errors"
	"os"
	"strconv"
)

type Errno uintptr

func (e Errno) Error() string {
	if 0 <= int(e) && int(e) < len(errors) {
		s := errors[e]
		if s != "" {
			return s
		}
	}
	return "errno " + strconv.Itoa(int(e))
}

func (e Errno) Is(target error) bool {
	switch target {
	case os.ErrPermission:
		return e == EACCES || e == EPERM
	case os.ErrExist:
		return e == EEXIST || e == ENOTEMPTY
	case os.ErrNotExist:
		return e == ENOENT
	case errorspkg.ErrUnsupported:
		return e == ENOSYS || e == ENOTSUP || e == EOPNOTSUPP
	}
	return false
}

func (e Errno) Temporary() bool {
	return e == EINTR || e == EMFILE || e == ENFILE || e.Timeout()
}

func (e Errno) Timeout() bool {
	return e == EAGAIN || e == EWOULDBLOCK || e == ETIMEDOUT
}

var (
	errEAGAIN error = EAGAIN
	errEINVAL error = EINVAL
	errENOENT error = ENOENT
)

func errnoErr(e Errno) error {
	switch e {
	case 0:
		return nil
	case EAGAIN:
		return errEAGAIN
	case EINVAL:
		return errEINVAL
	case ENOENT:
		return errENOENT
	}
	return e
}
