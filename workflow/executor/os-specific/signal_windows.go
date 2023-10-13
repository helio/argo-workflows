package os_specific

import (
	"fmt"
	argoerrors "github.com/argoproj/argo-workflows/v3/util/errors"
	"os"
	"syscall"
)

var (
	Term = os.Interrupt
)

func CanIgnoreSignal(s os.Signal) bool {
	fmt.Println("SIGNAL IGNORE WIN?", s)
	return false
}

func Kill(pid int, s syscall.Signal) error {
	if pid < 0 {
		pid = -pid // // we cannot kill a negative process on windows
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(s)
}

func Setpgid(a *syscall.SysProcAttr) {
	// this does not exist on windows
}

func Wait(process *os.Process) error {
	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		return err
	}
	if err := process.Release(); err != nil {
		return err
	}
	const STILL_ACTIVE = 259
	for {
		var ec uint32
		err := syscall.GetExitCodeProcess(handle, &ec)
		if err != nil {
			return os.NewSyscallError("GetExitCodeProcess", err)
		}
		if ec == STILL_ACTIVE {
			continue
		}
		if ec != 0 {
			return argoerrors.NewExitErr(int(ec))
		}
		return nil
	}
}
