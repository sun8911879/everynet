package internet

import (
	"errors"
	"syscall"
)

func InternetSet() error {
	lib := syscall.MustLoadDLL("Wininet.dll")
	c := lib.MustFindProc("InternetSetOptionA")
	r, _, _ := c.Call(uintptr(0), uintptr(39), uintptr(0), uintptr(0))
	if r != 0 {
		return errors.New("WINAPI InternetSetOptionA ERROR!")
	}
	return nil
}
