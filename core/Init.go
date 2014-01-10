package core

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"
)

//网站 or IP列表
var DomanList = &Doman{Value: make(map[string]*Attribute)}

//远程服务器信息
var Addr string

//列表数据类型
type Doman struct {
	Value map[string]*Attribute
}

type Attribute struct {
	Port   map[string]bool
	Virtue bool
}

//gfwlist列表--处理后的
const (
	GFWLIST = "http://cloudspeed.sinaapp.com/list.txt"
	INFO    = "http://cloudspeed.sinaapp.com/server_addr.txt"
)

func Auto() error {
	err := gfwlist()
	if err != nil {
		return err
	}
	err = info()
	if err != nil {
		return err
	}
	return nil
}

//获取gfwlist 列表
func gfwlist() error {
	req, err := http.NewRequest("GET", GFWLIST, nil)
	if err != nil {
		return err
	}
	//申请客户端
	client := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(15 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*15)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	//请求服务端
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("Get data error!")
	}
	defer resp.Body.Close()
	buffio := bufio.NewReader(resp.Body) //读入缓存
	for {
		line, err := buffio.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil {
			break
		}
		damon := strings.Split(line, ":")
		//判断是否为注释或者空
		if len(damon[0]) < 2 || (damon[0][:1] == "#" || damon[0][:1] == "/") {
			continue
		}
		//去掉换行
		if len(damon) == 1 {
			damon[0] = damon[0][:len(damon[0])-1]
		}

		DomanList.Value[damon[0]] = &Attribute{Port: make(map[string]bool), Virtue: false}

		//判断是否有端口
		if len(damon) > 1 {
			DomanList.Value[damon[0]].Virtue = true
			//去掉最后换行
			damon[len(damon)-1] = damon[len(damon)-1][:len(damon[len(damon)-1])-1]
		}
		for _, value := range damon[1:] {
			DomanList.Value[damon[0]].Port[value] = true
		}

	}
	return nil
}

func info() error {
	req, err := http.NewRequest("GET", INFO, nil)
	if err != nil {
		return err
	}
	//申请客户端
	client := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(15 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*15)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	//请求服务端
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("Get server_info error!")
	}
	buffio := bufio.NewReader(resp.Body) //读入缓存
	line, err := buffio.ReadString('\n') //以'\n'为结束符读入一行
	if err != nil {
		return errors.New("Get server_info error!")
	}
	Addr = line[:len(line)-1]
	return nil
}
