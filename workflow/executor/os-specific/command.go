package os_specific

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
)

var logger = log.WithField("argo", true)

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
		_ = cmd.Wait()
	}

	return closer, nil
}

func isTerminal(stdin io.Reader) bool {
	f, ok := stdin.(*os.File)
	return ok && term.IsTerminal(int(f.Fd()))
}
