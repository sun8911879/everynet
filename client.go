package main

import (
	"fmt"
	"github.com/sun8911879/everynet/core"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//监听退出
	go core.Os_kill()
	//加载数据
	err := core.Auto()
	if nil != err {
		fmt.Println(err)
		os.Exit(0)
	}
	core.Handle()
	return
}
