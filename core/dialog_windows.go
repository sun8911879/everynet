package core

import (
	"encoding/gob"
	"github.com/sun8911879/everynet/tools/memory"
	"github.com/sun8911879/everynet/tools/tcp"
	"net"
	"unsafe"
)

func (Tcp *Request) Secret() {
	//建立连接
	if Tcp.Remote == nil {
		//建立远程连接
		Tcp.Remote, Tcp.err = net.Dial("tcp", Addr)
		if Tcp.err != nil {
			return
		}
		//写入信息
		methods := make([]byte, 99)
		methods = []byte(Tcp.Addr)
		for i := len(methods); i < 99; i++ {
			methods = append(methods, 0)
		}
		Tcp.Remote.Write(methods)
		Tcp.enc = gob.NewEncoder(Tcp.Remote)
		//读取数据
		go tcp.GobRead(Tcp.Accept, Tcp.Remote)
	}
	//写入头
	Tcp.err = Tcp.GobWriter()
	return
}

func (Tcp *Request) Secrets() {
	//建立远程连接
	Tcp.Remote, Tcp.err = net.Dial("tcp", Addr)
	//返回数据
	if Tcp.err != nil {
		return
	}
	defer Tcp.Remote.Close()
	//写入信息
	methods := make([]byte, 99)
	methods = []byte(Tcp.Addr)
	for i := len(methods); i < 99; i++ {
		methods = append(methods, 0)
	}

	Tcp.Remote.Write(methods)
	//写入数据
	go tcp.GobWriter(Tcp.Remote, Tcp.Accept)
	tcp.GobRead(Tcp.Accept, Tcp.Remote)
	return
}

//编码写入
func (Tcp *Request) GobWriter() error {
	err := Tcp.enc.Encode([]byte(Tcp.Head))
	//如果POST 写入数据
	if Tcp.Pact == "POST" {
		Tcp.PostCopy()
	}
	return err
}

func (Tcp *Request) GobPostCopy() (n int) {
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
			ew := Tcp.enc.Encode(buf[:nr])
			if nr > 0 {
				n += int(nr)
			}
			if ew != nil {
				Tcp.err = ew
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
