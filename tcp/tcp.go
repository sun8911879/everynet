package tcp

import (
	"errors"
	"github.com/sun8911879/everynet/memory"
	"unsafe"
)

var ErrShortWrite = errors.New("short write")
var EOF = errors.New("EOF")

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type WinReader interface {
	Read(p []byte) (n int, err error)
	Close() error
}

type WinWriter interface {
	Write(p []byte) (n int, err error)
	Close() error
}

func Copy(dst Writer, src Reader) (n int, err error) {
	alloc := memory.Alloc(uintptr(32 * 1024))
	buf := (*[1 << 30]byte)(unsafe.Pointer(alloc))[:32*1024]
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				n += int(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		} else {
			break
		}
		if er == EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	memory.Free(alloc, uintptr(32*1024))
	alloc = nil
	buf = nil
	return n, err
}

//POST拷贝
func PostCopy(dst Writer, src Reader) (n int, err error) {
	alloc := memory.Alloc(uintptr(32 * 1024))
	buf := (*[1 << 30]byte)(unsafe.Pointer(alloc))[:32*1024]
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				n += int(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er == EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	memory.Free(alloc, uintptr(32*1024))
	alloc = nil
	buf = nil
	return n, err
}

//拷贝完成关闭连接
func WinCopy(dst WinWriter, src WinReader) (n int, err error) {
	alloc := memory.Alloc(uintptr(32 * 1024))
	buf := (*[1 << 30]byte)(unsafe.Pointer(alloc))[:32*1024]
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				n += int(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		} else {
			break
		}
		if er == EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	src.Close()
	dst.Close()
	memory.Free(alloc, uintptr(32*1024))
	alloc = nil
	buf = nil
	return n, err
}
