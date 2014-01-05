package core

import (
	"os"
	"os/signal"
	"syscall"
)

func Os_kill() {
	schan := make(chan os.Signal)
	go signal.Notify(schan, syscall.SIGINT, syscall.SIGTERM)
	<-schan
	Net_Off()
	os.Exit(0)
}
