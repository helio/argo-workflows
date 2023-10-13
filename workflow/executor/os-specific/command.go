package os_specific

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
	"io"
	"os"
)

var logger = log.WithField("argo", true)

func isTerminal(stdin io.Reader) bool {
	f, ok := stdin.(*os.File)
	return ok && term.IsTerminal(int(f.Fd()))
}
