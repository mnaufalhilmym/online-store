package logger

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

var elog = log.New(os.Stderr, "[ERROR]", 1)

func (l *logger) Error(message interface{}, options ...*Options) {
	msg := []interface{}{"[" + l.prefix + "]", message}

	if options == nil || options[0].IsPrintStack {
		msg = append(msg, fmt.Sprintf("\n%s", debug.Stack()))
	}

	elog.Println(msg...)

	if len(options) > 0 && options[0].IsExit {
		exitCode := 1
		if options[0].ExitCode > 1 {
			exitCode = options[0].ExitCode
		}

		os.Exit(exitCode)
	}
}
