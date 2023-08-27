package logger

import (
	"log"
	"os"
)

var plog = log.New(os.Stderr, "[PANIC]", 2)

func (l *logger) Panic(message interface{}, options ...Options) {
	msg := []interface{}{"[" + l.prefix + "]", message}

	plog.Panicln(msg...)
}
