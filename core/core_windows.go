package core

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/sun8911879/everynet/tcp"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var domon_regexp, _ = regexp.Compile(`[^.]+\.(com|cn|net|org|edu|gov|biz|tv|me|pro|name|cc|co|info|cm)(\.(cn|us|hk|tw|uk|it|fr|br|in|de))?`)

var HTTP_200 = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")

var HTTP_407 = []byte("HTTP/1.1 407 Unauthorized\r\n\r\n")

const (
	CONNECT = "CONNECT"
	HTTP    = " HTTP"
)

//进行处理
func Handle() {
	//监听
	go Listen()
	//开启网络设置
	Net_On()
	//绑定事件
	err := Minimize()
	if err != nil {
		fmt.Println(err)
	}
	//安全退出
	Net_Off()
	//终止执行
	os.Exit(0)
}

//通信存储类型
type Request struct {
	Host   string        //域名
	Source string        //源域名
	Path   string        //路径
	GET    string        //请求地址
	Pact   string        //协议 GET POST HTTPS
	Port   string        //端口
	Addr   string        //tcp远程地址
	Cache  string        //缓存
	Key    string        //key--判断是否需要代理
	Length int           //Body长度
	chanl  bool          //判断通道是否已经打开
	err    error         //错误类型
	Remote net.Conn      //远端TCP请求
	Accept net.Conn      //本地TCP请求
	bufio  *bufio.Reader //读取缓存bufio
	enc    *gob.Encoder  //gob对象
}

//监听
func Listen() {
	//监听端口
	ln, err := net.Listen("tcp", ":5316")
	if nil != err {
		os.Exit(0)
	}
	//处理数据
	for {
		Tcp := &Request{Port: "80"}
		Tcp.Accept, err = ln.Accept()
		if nil != err {
			continue
		}
		go Tcp.Serve()
	}
}

//进行服务
func (Tcp *Request) Serve() {
	//defer Tcp.Accept.Close()
	Tcp.bufio = bufio.NewReader(Tcp.Accept) //读入缓存
	//判断协议
	line, err := Tcp.bufio.ReadString('\n')
	pact_index := strings.Index(line, " ")
	if err != nil || pact_index < 3 {
		return
	}
	//协议
	Tcp.Pact = line[:pact_index]
	//获取路径
	path_index := strings.Index(line, HTTP)
	if path_index == -1 || pact_index >= path_index {
		return
	}
	//路径
	Tcp.Path = line[pact_index+1 : path_index]
	//区分处理数据流
	switch Tcp.Pact {
	case CONNECT:
		Tcp.HTTPS()
		break
	default:
		Tcp.Cache = line
		Tcp.HTTP()
		break
	}
}

//HTTP协议处理
func (Tcp *Request) HTTP() {
	//关闭连接
	defer func() {
		if Tcp.Remote != nil {
			Tcp.Remote.Close()
		}
	}()
	for {
		if Tcp.chanl {
			line, err := Tcp.bufio.ReadString('\n')
			pact_index := strings.Index(line, " ")
			if err != nil || pact_index < 3 {
				return
			}
			//协议
			Tcp.Pact = line[:pact_index]
			//获取路径
			path_index := strings.Index(line, HTTP)
			if path_index == -1 || pact_index >= path_index {
				return
			}
			//路径
			Tcp.Path = line[pact_index+1 : path_index]
			Tcp.Cache = line
		}
		if Tcp.err != nil {
			return
		}
		Tcp.Headr()
		Tcp.Wall()
		Tcp.Ship()
	}
}

//HTTP协议头部处理
func (Tcp *Request) Headr() {
	//替换地址
	Tcp.GET = strings.Replace(Tcp.Path, "http://", "", 1)
	GET_Index := strings.Index(Tcp.GET, "/")
	if GET_Index == -1 {
		Tcp.err = errors.New("GET PATH ERROR!")
		return
	}
	Tcp.Source = Tcp.GET[:GET_Index]
	Tcp.Host = ""
	Tcp.Video()
	if Tcp.Host == "" {
		Tcp.GET = Tcp.GET[GET_Index:]
	}
	//替换源地址
	Tcp.Cache = strings.Replace(Tcp.Cache, Tcp.Path, Tcp.GET, 1)
	//判断域名进行申请
	for {
		line, err := Tcp.bufio.ReadString('\n')
		if err != nil {
			Tcp.err = err
			break
		}
		if len(line) <= 2 {
			Tcp.Cache = Tcp.Cache + line
			break
		}
		//获取域名
		Host := strings.Index(line, "Host:")
		if Host == 0 && len(line) > 8 {
			//判断域名--是否需要更改(判断是否被广告替换掉)
			if Tcp.Host != "" {
				line = "Host: " + Tcp.Host + "\r\n"
			} else {
				Tcp.Host = line[6 : len(line)-2]
			}
		}
		//Content-Length
		if strings.Index(line, "Content-Length:") != -1 {
			Tcp.Length, _ = strconv.Atoi(line[16 : len(line)-2])
		}
		//更改代理标识符
		if strings.Index(line, "Proxy-Connection:") != -1 {
			line = strings.Replace(line, "Proxy-Connection", "Connection", 1)
		}
		Tcp.Cache = Tcp.Cache + line
	}
	//判断端口
	port_index := strings.LastIndex(Tcp.Host, ":")
	if port_index != -1 && len(Tcp.Host) > port_index+1 {
		Tcp.Port = Tcp.Host[port_index+1:]
	}
	//判断通道地址是否一样.不一样关闭远程连接.重置
	if Tcp.chanl == true && Tcp.Addr != Tcp.Host+":"+Tcp.Port {
		Tcp.Remote.Close()
		Tcp.Remote = nil
		Tcp.chanl = false
	}
	Tcp.Addr = Tcp.Host + ":" + Tcp.Port
}

//HTTP传输数据
func (Tcp *Request) Ship() {
	//加密传输
	if Tcp.Key != "" {
		Tcp.Secret()
		return
	}
	//建立连接
	if Tcp.chanl == false {
		Tcp.Remote, Tcp.err = net.Dial("tcp", Tcp.Addr)
		if Tcp.err != nil {
			return
		}
	}
	//写入数据
	Tcp.Remote.Write([]byte(Tcp.Cache))
	//如果POST 写入数据
	if Tcp.Pact == "POST" {
		if Tcp.Length == 0 {
			Tcp.err = errors.New("POST Length Too Short")
			return
		}
		tcp.PostCopy(Tcp.Remote, Tcp.bufio, Tcp.Length)
	}
	Tcp.Cache = ""
	//读取数据--判断通道是否已经开启
	if Tcp.chanl == false {
		go tcp.WinCopy(Tcp.Accept, Tcp.Remote)
		Tcp.chanl = true
	}
}

//HTTPS协议处理
func (Tcp *Request) HTTPS() {
	//去掉无用信息
	for {
		line, err := Tcp.bufio.ReadString('\n')
		if len(line) <= 2 {
			break
		}
		if err != nil {
			return
		}
	}
	//设置地址
	Tcp.Addr = Tcp.Path
	//判断远程地址
	host_index := strings.LastIndex(Tcp.Path, ":")
	if host_index != -1 && len(Tcp.Path) > host_index+1 {
		Tcp.Host = Tcp.Path[:host_index]
		Tcp.Port = Tcp.Path[host_index+1:]
	}
	//返回协议
	Tcp.Accept.Write(HTTP_200)
	Tcp.Ships()
}

//HTTPS传输数据
func (Tcp *Request) Ships() {
	Tcp.Wall()
	//加密传输
	if Tcp.Key != "" {
		Tcp.Secrets()
		return
	}
	//建立连接
	Tcp.Remote, Tcp.err = net.Dial("tcp", Tcp.Addr)
	if Tcp.err != nil {
		return
	}
	defer Tcp.Remote.Close()
	//写入数据
	go tcp.Copy(Tcp.Remote, Tcp.bufio)
	//读取数据
	tcp.Copy(Tcp.Accept, Tcp.Remote)
}

//根据列表.进行过滤是否需要加密传输--支持HTTP-HTTPS
func (Tcp *Request) Wall() {
	if DomanList.Value[Tcp.Host] != nil {
		Tcp.Key = Tcp.Host
	}
	if Tcp.Key == "" {
		damon_regexp := domon_regexp.FindStringSubmatch(Tcp.Host)
		if len(damon_regexp) > 0 {
			damon := "." + damon_regexp[0]
			if DomanList.Value[damon] != nil {
				Tcp.Key = damon
			}
		}
	}
	//判断是否需要加密传输
	if Tcp.Key != "" {
		//判断端口还是全部代理
		//全部代理
		if DomanList.Value[Tcp.Key].Virtue == false {
			return
		}
		//端口正确
		if DomanList.Value[Tcp.Key].Port[Tcp.Port] == true {
			return
		}
		Tcp.Key = ""
	}
}
