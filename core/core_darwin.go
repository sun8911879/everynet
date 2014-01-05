package core

import (
	"errors"
	"fmt"
	"github.com/sun8911879/everynet/tcp"
	"net"
	"os"
	"regexp"
)

var domon_regexp, _ = regexp.Compile(`[^.]+\.(com|cn|net|org|edu|gov|biz|tv|me|pro|name|cc|co|info|cm)(\.(cn|us|hk|tw|uk|it|fr|br|in|de))?`)

//进行处理
func Handle() {
	//监听端口
	ln, err := net.Listen("tcp", ":5316")
	if nil != err {
		os.Exit(0)
	}
	//开启网络设置
	Net_On()
	//处理数据
	for {
		conn, err := ln.Accept()
		if nil != err {
			continue
		}
		go Handle_Flow(conn)
	}
}

//协议
type Protocol struct {
	Conn net.Conn
	Req  request
}

//处理数据流
func Handle_Flow(conn net.Conn) {
	defer conn.Close()
	core := &Protocol{Conn: conn}
	//版本信息校验
	err := core.check()
	if err != nil {
		return
	}
	//返回协议
	err = core.back()
	if err != nil {
		return
	}
	//读取资源
	err = core.req()
	if nil != err {
		return
	}

	if DomanList.Value[core.Req.remote] != nil {
		core.Req.key = core.Req.remote
	}

	if core.Req.key == "" {
		damon_regexp := domon_regexp.FindStringSubmatch(core.Req.remote)
		if len(damon_regexp) > 0 {
			damon := "." + damon_regexp[0]
			if DomanList.Value[damon] != nil {
				core.Req.key = damon
			}
		}
	}

	//判断是否需要加密传输
	if core.Req.key != "" {
		//判断端口是否需要全部代理
		//全部代理
		if DomanList.Value[core.Req.key].Virtue == false {
			core.remote()
			return
		}
		//端口正确
		if DomanList.Value[core.Req.key].Port[core.Req.prot] == true {
			core.remote()
			return
		}
	}

	//创建连接
	remote, err := net.Dial(core.Req.reqtype, core.Req.addr)
	//返回数据
	if nil != err {
		core.gen(4)
		return
	}
	core.gen(0)
	go tcp.Copy(remote, core.Conn)
	tcp.Copy(core.Conn, remote)
	remote.Close()
	return
}

//头部协议
type Pact struct {
	ver      [1]uint8
	nmethods [1]uint8
	methods  [99]uint8
}

//版本信息校验
func (protocol *Protocol) check() (err error) {
	pact := &Pact{}
	_, err = protocol.Conn.Read(pact.ver[:])
	if nil != err {
		return errors.New("check error")
	}
	_, err = protocol.Conn.Read(pact.nmethods[:])
	if nil != err {
		return errors.New("check error")
	}
	_, err = protocol.Conn.Read(pact.methods[:int(pact.nmethods[0])])
	if nil != err {
		return errors.New("check error")
	}
	return nil
}

//返回cocks响应协议--是否接收这次请求
func (protocol *Protocol) back() (err error) {
	buff := make([]uint8, 2)
	buff[0] = 5
	buff[1] = 0
	_, err = protocol.Conn.Write(buff)
	if nil != err {
		return errors.New("back error")
	}
	return nil
}

//客户端请求资源
type request struct {
	ver       uint8     // 	版本信息:socks v5
	cmd       uint8     // 	连接信息:CONNECT: 0x01, BIND:0x02, UDP ASSOCIATE: 0x03
	rsv       uint8     //	版权?
	atyp      uint8     //	连接信息:IPv4 or IPv6 or 域名
	dst_addr  [256]byte //	远程地址
	dst_port  [2]uint8  //  远程端口--uint8未转码
	dst_port2 uint16    //	远程端口--处理后.uint8转码后
	reqtype   string    //	socket:tcp or udp
	//程序需要--不是socks协议
	prot   string //	端口--字符串一次转码
	addr   string //	远程地址--包含端口
	key    string //	map--key支持
	remote string //  	过滤地址--支持域名和IP
}

//请求资源处理
func (protocol *Protocol) req() (err error) {
	buff := make([]byte, 4)
	_, err = protocol.Conn.Read(buff)
	if nil != err {
		return err
	}
	protocol.Req.ver, protocol.Req.cmd, protocol.Req.rsv, protocol.Req.atyp = buff[0], buff[1], buff[2], buff[3]

	if 5 != protocol.Req.ver || 0 != protocol.Req.rsv {
		return errors.New("Request Message VER or RSV error!")
	}
	switch protocol.Req.atyp {
	case 1: //ip v4
		_, err = protocol.Conn.Read(protocol.Req.dst_addr[:4])
	case 4: //ipv6
		_, err = protocol.Conn.Read(protocol.Req.dst_addr[:16])
	case 3: //DOMANNAME
		_, err = protocol.Conn.Read(protocol.Req.dst_addr[:1])
		_, err = protocol.Conn.Read(protocol.Req.dst_addr[1 : 1+int(protocol.Req.dst_addr[0])])
	}
	if nil != err {
		return errors.New("Request IP error!")
	}
	_, err = protocol.Conn.Read(protocol.Req.dst_port[:2])
	if nil != err {
		return errors.New("Request PROT error!")
	}
	//地址
	protocol.Req.remote = string(protocol.Req.dst_addr[1 : 1+protocol.Req.dst_addr[0]])
	//获取端口
	protocol.Req.dst_port2 = uint16(uint16(protocol.Req.dst_port[0])*256 + uint16(protocol.Req.dst_port[1]))
	//判断类型 tcp or udp
	switch protocol.Req.cmd {
	case 1:
		protocol.Req.reqtype = "tcp"
	case 3:
		protocol.Req.reqtype = "udp"
	}
	//获取远程地址
	switch protocol.Req.atyp {
	case 1: // ipv4
		protocol.Req.addr = fmt.Sprintf("%d,%d,%d,%d:%d", protocol.Req.dst_addr[0], protocol.Req.dst_addr[1], protocol.Req.dst_addr[2], protocol.Req.dst_addr[3], protocol.Req.dst_port2)
		break
	case 3: //DOMANNAME
		protocol.Req.prot = fmt.Sprintf(":%d", protocol.Req.dst_port2)
		protocol.Req.addr = protocol.Req.remote + protocol.Req.prot
		break
		//case 4: //ipv6
	}
	return
}

//返回数据协议
type backing struct {
	ver  uint8
	rep  uint8
	rsv  uint8
	atyp uint8
	buf  [10]uint8
}

//返回数据协议
func (protocol *Protocol) gen(rep uint8) {
	back := &backing{
		ver:  5,
		rep:  rep,
		rsv:  0,
		atyp: 1,
	}
	back.buf[0], back.buf[1], back.buf[2], back.buf[3] = back.ver, back.rep, back.rsv, back.atyp
	protocol.Conn.Write(back.buf[:10])
	return
}
