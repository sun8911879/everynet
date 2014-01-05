package main

import (
	"fmt"
	"github.com/sun8911879/everynet/tcp"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":5317")
	if nil != err {
		fmt.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		if nil != err {
			fmt.Println("Accept Error!")
			continue
		}
		go Handle(conn)
	}
}

func Handle(conn net.Conn) {
	defer conn.Close()
	//读取99长度协议
	methods := make([]byte, 99)
	conn.Read(methods)
	remote, err := net.Dial("tcp", string(methods))
	if err != nil {
		return
	}
	defer remote.Close()
	go tcp.GobRead(remote, conn)
	tcp.GobWriter(conn, remote)
	return
}
