package log

import (
	"runtime/debug"
)

func SaveLogService(location string, message string, printStack bool) error {
	stack := new(string)
	if printStack {
		_stack := string(debug.Stack())
		stack = &_stack
	}
	if _, err := LogRepository().Create(&logModel{
		Location: &location,
		Message:  &message,
		Stack:    stack,
	}); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
