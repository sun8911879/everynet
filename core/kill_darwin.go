package core

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
)

type signalHandler func(s os.Signal, arg interface{})

type signalSet struct {
	m map[os.Signal]signalHandler
}

func signalSetNew() *signalSet {
	ss := new(signalSet)
	ss.m = make(map[os.Signal]signalHandler)
	return ss
}

func (set *signalSet) register(s os.Signal, handler signalHandler) {
	if _, found := set.m[s]; !found {
		set.m[s] = handler
	}
}

func (set *signalSet) handle(sig os.Signal, arg interface{}) (err error) {
	if _, found := set.m[sig]; found {
		set.m[sig](sig, arg)
		return nil
	} else {
		if sig.String() == "terminated" {
			return errors.New("os normal exit")
		}
	}
	return nil
}

func Os_kill() {
	ss := signalSetNew()
	handler := func(s os.Signal, arg interface{}) {
		Net_Off()
		os.Exit(0)
	}
	ss.register(syscall.SIGINT, handler)
	ss.register(syscall.SIGUSR1, handler)
	ss.register(syscall.SIGUSR2, handler)
	for {
		c := make(chan os.Signal)
		var sigs []os.Signal
		for sig := range ss.m {
			sigs = append(sigs, sig)
		}
		signal.Notify(c)
		sig := <-c
		err := ss.handle(sig, nil)
		if err != nil {
			Net_Off()
			os.Exit(0)
		}
	}
}
