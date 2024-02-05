package log

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type MyFormatter struct {
	PID int
}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05,000")
	var newLog string
	newLog = fmt.Sprintf("[PID:%d %s %s]: %s\n", m.PID, timestamp, MyLevel{Level: entry.Level}, entry.Message)

	b.WriteString(newLog)
	return b.Bytes(), nil
}

type MyLevel struct {
	logrus.Level
}

func (level MyLevel) String() string {
	if b, err := level.MarshalText(); err == nil {
		return strings.ToUpper(string(b))
	} else {
		return "unknown"
	}
}
