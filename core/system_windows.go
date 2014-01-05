package core

import (
	"github.com/lxn/walk"
	"os/exec"
)

//最小化窗口
func Minimize() error {
	mw, err := walk.NewMainWindow()
	if err != nil {
		return err
	}
	//加载图片文件目录
	icon, err := walk.NewIconFromFile("./Resources/logo.ico")
	if err != nil {
		return err
	}
	ni, err := walk.NewNotifyIcon()
	if err != nil {
		return err
	}
	if err := ni.SetIcon(icon); err != nil {
		return err
	}
	//鼠标放上信息
	if err := ni.SetToolTip("everynet加速器,点击设置或退出."); err != nil {
		return err
	}
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		//everynet
		if err := ni.ShowCustom(
			//点击显示
			"everynet加速器",
			"everynet正在为您的网络进行加速中...."); err != nil {
			return
		}
	})
	// 菜单目录设置
	exitAction := walk.NewAction()
	if err := exitAction.SetText("退出加速器"); err != nil {
		return err
	}
	exitAction.Triggered().Attach(func() {
		walk.App().Exit(0)
	})
	// 标题设置
	titleAction := walk.NewAction()
	if err := titleAction.SetText("万象加速器1.0"); err != nil {
		return err
	}
	//添加按钮事件
	if err := ni.ContextMenu().Actions().Add(titleAction); err != nil {
		return err
	}
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		return err
	}
	// 设置通知图标可见
	if err := ni.SetVisible(true); err != nil {
		return err
	}
	// 默认显示
	if err := ni.ShowInfo("everynet已经启动", "右键图标进行网络加速设置"); err != nil {
		return err
	}
	// 开始运行
	mw.Run()
	defer ni.Dispose()
	return nil
}

//设置系统--开启
func Net_On() {
	networksetup_reg_add := exec.Command(`reg`, `add`, `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, `/v`, `ProxyServer`, `/t`, `REG_SZ`, `/d`, `http=127.0.0.1:5316;https=127.0.0.1:5316`, `/f`)
	networksetup_reg_add.Output()
	networksetup_reg_on := exec.Command(`reg`, `add`, `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, `/v`, `ProxyEnable`, `/t`, `REG_DWORD`, `/d`, `1`, `/f`)
	networksetup_reg_on.Output()
}

//设置系统网络代理--关闭
func Net_Off() {
	networksetup_reg_on := exec.Command(`reg`, `add`, `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, `/v`, `ProxyEnable`, `/t`, `REG_DWORD`, `/d`, `0`, `/f`)
	networksetup_reg_on.Output()
}
