package os_specific

import (
	"fmt"
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
	//fmt.Println("#### process wait")
	//stat, err := process.Wait()
	//fmt.Println(">>>>>>> end wait")
	//if stat.ExitCode() != 0 {
	//	return errors.NewExitErr(stat.ExitCode())
	//}
	//return err
	return nil
}
