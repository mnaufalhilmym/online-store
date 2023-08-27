package logger

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

var nlog = log.New(os.Stdout, "[LOG]", 0)

func (l *logger) Log(message interface{}, options ...*Options) {
	msg := []interface{}{"[" + l.prefix + "]", message}

	if len(options) > 0 && options[0].IsPrintStack {
		msg = append(msg, fmt.Sprintf("\n%s", debug.Stack()))
	}

	nlog.Println(msg...)
}
