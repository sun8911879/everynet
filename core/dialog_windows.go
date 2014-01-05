package core

import (
	"encoding/gob"
	"github.com/sun8911879/everynet/tcp"
	"net"
)

func (Tcp *Request) Secret() {
	//建立连接
	if Tcp.chanl == false {
		//建立远程连接
		Tcp.Remote, Tcp.err = net.Dial("tcp", Addr)
		if Tcp.err != nil {
			return
		}
	}
	if Tcp.chanl == false {
		//写入信息
		methods := make([]byte, 99)
		methods = []byte(Tcp.Addr)
		for i := len(methods); i < 99; i++ {
			methods = append(methods, 0)
		}
		Tcp.Remote.Write(methods)
		Tcp.enc = gob.NewEncoder(Tcp.Remote)
	}
	//写入头
	Tcp.err = Tcp.GobHeadWriter()
	if Tcp.err != nil {
		return
	}
	//读取数据--判断通道是否已经开启
	if Tcp.chanl == false {
		go tcp.GobRead(Tcp.Accept, Tcp.Remote)
		Tcp.chanl = true
	}
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
func (Tcp *Request) GobHeadWriter() error {
	err := Tcp.enc.Encode([]byte(Tcp.Cache))
	return err
}
