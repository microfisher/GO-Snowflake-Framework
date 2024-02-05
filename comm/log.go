package comm

import (
	"fmt"
	"snowflake/log"
)

// Debug .
func Debug(identifer, format string, args ...any) {
	var message = fmt.Sprintf(format, args...)
	log.Debugf("%s -> %s", identifer, message)
}

// Info .
func Info(identifer, format string, args ...any) {
	var message = fmt.Sprintf(format, args...)
	log.Infof("%s -> %s", identifer, message)
}

// Warn .
func Warn(identifer, format string, args ...any) {
	var message = fmt.Sprintf(format, args...)
	log.Warnf("%s -> %s", identifer, message)
}

// Error .
func Error(identifer, format string, args ...any) {
	var message = fmt.Sprintf(format, args...)
	log.Errorf("%s %s -> %s", identifer, message, EventWarn)
}
