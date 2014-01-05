package core

import (
	"os/exec"
	"syscall"
)

//设置osx系统--开启
func Net_On() {
	//文件句柄
	Linux_opens()
	//WI-FI
	networksetup_socks_wifi := exec.Command(`networksetup`, `-setsocksfirewallproxy`, `Wi-Fi`, `127.0.0.1`, `5316`)
	networksetup_socks_wifi.Output()
	networksetup_socks_wifi_switch := exec.Command(`networksetup`, `-setsocksfirewallproxystate`, `Wi-Fi`, `on`)
	networksetup_socks_wifi_switch.Output()
	//Ethernet
	networksetup_socks_ethernet := exec.Command(`networksetup`, `-setsocksfirewallproxy`, `Ethernet`, `127.0.0.1`, `5316`)
	networksetup_socks_ethernet.Output()
	networksetup_socks_ethernet_switch := exec.Command(`networksetup`, `-setsocksfirewallproxystate`, `Ethernet`, `on`)
	networksetup_socks_ethernet_switch.Output()
}

//unix文件句柄
func Linux_opens() error {
	limit := &syscall.Rlimit{
		Cur: 100000,
		Max: 100000,
	}
	err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, limit)
	return err
}

//设置系统网络代理--关闭
func Net_Off() {
	//WI-FI
	networksetup_socks_wifi_switch := exec.Command(`networksetup`, `-setsocksfirewallproxystate`, `Wi-Fi`, `off`)
	networksetup_socks_wifi_switch.Output()
	//Ethernet
	networksetup_socks_ethernet_switch := exec.Command(`networksetup`, `-setsocksfirewallproxystate`, `Ethernet`, `off`)
	networksetup_socks_ethernet_switch.Output()
}
