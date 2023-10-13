package os_specific

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func StartCommand(cmd *exec.Cmd) (func(), error) {
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}

	if isTerminal(cmd.Stdin) {
		logger.Warn("TTY detected but is not supported on windows")
	}
	return simpleStart(cmd)
}

func simpleStart(cmd *exec.Cmd) (func(), error) {
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	closer := func() {
		cmd.Cancel = func() error {
			fmt.Println("#######")
			return nil
		}
		cmd.WaitDelay = 100 * time.Millisecond
		fmt.Println("Starting to wait")
		if err := cmd.Wait(); err != nil {
			fmt.Println("!!!!!!!!!!", err)
		}
	}

	return closer, nil
}
