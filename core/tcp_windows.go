package core

import (
	"errors"
	"github.com/sun8911879/everynet/tools/memory"
	"unsafe"
)

var ErrShortWrite = errors.New("short write")
var EOF = errors.New("EOF")
var POST = errors.New("POST Length Too Short")

//POST拷贝
func (Tcp *Request) PostCopy() (n int) {
	if Tcp.Length < 1 {
		Tcp.err = POST
		return n
	}
	if Tcp.Length <= 1024 {
		Tcp.alloc = 1024
	} else {
		Tcp.alloc = 32 * 1024
	}

	alloc := memory.Alloc(uintptr(Tcp.alloc))
	buf := (*[1 << 30]byte)(unsafe.Pointer(alloc))[:Tcp.alloc]
	for {
		nr, er := Tcp.src.Read(buf)
		if nr > 0 {
			nw, ew := Tcp.Remote.Write(buf[0:nr])
			if nw > 0 {
				n += int(nw)
			}
			if ew != nil {
				Tcp.err = ew
				break
			}
			if nr != nw {
				Tcp.err = ErrShortWrite
				break
			}
		}
		if n >= Tcp.Length {
			break
		}
		if er == EOF {
			break
		}
		if er != nil {
			Tcp.err = er
			break
		}
	}
	memory.Free(alloc, uintptr(Tcp.alloc))
	alloc = nil
	buf = nil
	return n
}

//拷贝完成关闭连接
func (Tcp *Request) WinCopy() (n int) {
	alloc := memory.Alloc(uintptr(32 * 1024))
	buf := (*[1 << 30]byte)(unsafe.Pointer(alloc))[:32*1024]
	for {
		if Tcp.Remote == nil {
			return 0
		}
		nr, er := Tcp.Remote.Read(buf)
		if nr > 0 {
			nw, ew := Tcp.Accept.Write(buf[0:nr])
			if nw > 0 {
				n += int(nw)
			}
			if ew != nil {
				Tcp.err = ew
				break
			}
			if nr != nw {
				Tcp.err = ErrShortWrite
				break
			}
		}
		if er == EOF {
			break
		}
		if er != nil {
			Tcp.err = er
			break
		}
	}
	if Tcp.Remote != nil {
		Tcp.Remote.Close()
	}
	Tcp.Accept.Close()
	memory.Free(alloc, uintptr(32*1024))
	alloc = nil
	buf = nil
	return n
}
