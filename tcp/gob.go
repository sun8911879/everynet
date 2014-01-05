package tcp

import (
	"encoding/gob"
	"github.com/sun8911879/everynet/memory"
	"unsafe"
)

type Encoder struct {
	Dst   *gob.Encoder
	Cache []byte
}

//解码读
func GobRead(dst Writer, src Reader) (n int, err error) {
	alloc := memory.Alloc(uintptr(32 * 1024))
	buf := (*[1 << 30]byte)(unsafe.Pointer(alloc))[:32*1024]
	dec := gob.NewDecoder(src)
	for {
		er := dec.Decode(&buf)
		if len(buf) < 1 {
			break
		}
		if er != nil {
			err = er
			break
		}
		nw, ew := dst.Write(buf[:len(buf)])
		if nw > 0 {
			n += int(nw)
		}
		if ew != nil {
			err = ew
			break
		}
	}
	memory.Free(alloc, uintptr(32*1024))
	alloc = nil
	buf = nil
	return n, err
}

//编码写入
func GobWriter(dst Writer, src Reader) (n int, err error) {
	alloc := memory.Alloc(uintptr(32 * 1024))
	buf := (*[1 << 30]byte)(unsafe.Pointer(alloc))[:32*1024]
	enc := gob.NewEncoder(dst)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			ew := enc.Encode(buf[:nr])
			if nr > 0 {
				n += int(nr)
			}
			if ew != nil {
				err = ew
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
