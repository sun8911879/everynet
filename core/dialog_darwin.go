package core

import (
	"github.com/sun8911879/everynet/tcp"
	"net"
)

func (protocol *Protocol) remote() {
	defer protocol.Conn.Close()
	//建立远程连接
	remote, err := net.Dial("tcp", Addr)
	//返回数据
	if nil != err {
		protocol.gen(4)
		return
	}
	defer remote.Close()
	protocol.gen(0)
	//写入信息
	methods := make([]byte, 99)
	methods = []byte(protocol.Req.addr)
	for i := len(methods); i < 99; i++ {
		methods = append(methods, 0)
	}
	remote.Write(methods)
	go tcp.GobWriter(remote, protocol.Conn)
	tcp.GobRead(protocol.Conn, remote)
	return
}
