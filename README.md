# everynet 网络加速器
golang编写 SOCKS5 和 http,https代理

golang自有 gob编码通信.

并发编程.多TCP同时连接服务端(支持浏览器TCP复用).速度是SSH等几倍

# OS X
OS X下 SOCKS5代理

通过更改networksetup实现

支持WIFI和Ethernet

无协议分析.纯代理.简单快速

# windows
windows下 HTTP,HTTPS代理

有协议分析 HTTP,HTTPS代理 支持更改HTTP请求

通过更改注册表(IE代理)实现

注册表改动: HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings

实现替换优酷播放器swf源.去广告目的

后期可以屏蔽广告联盟js等广告

windows下或许有少许bug.基本使用无问题

1.目前不完善.HTTP POST数据请求body部分.只能写完(具体效果看浏览器-有时候抽风会等待)

2.注册表改后无法实时刷新IE更改到代理.打算用cgo-wininet-InternetSetOption来实现刷新.目前没搞定求大神

## 交流

新浪微博：[雪虎](http://weibo.com/sun8911879)

##注释
由于初期时间紧.代码略烂.请见谅(此项目不保证更新)

##安装
安装方法：

	go get github.com/sun8911879/everynet

windows安装：

	go build

切记不要 go build client.go

OS X安装：
	
	go build client.go

服务端：

	go build server.go