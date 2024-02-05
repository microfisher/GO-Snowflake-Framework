package log

import (
	"io"
	"os"
	"strings"

	lj "github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type Level uint32

const (
	// 与logrus.InfoLevel 定义保持一至
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	traceLevel
)

// NewKLLogger 创建日志对象
func NewKLLogger(level Level, filename string, withSTDOUT bool) *KLLogger {
	var outputs []io.Writer
	if withSTDOUT || filename == "" || strings.ToUpper(filename) == "STDOUT" {
		outputs = append(outputs, os.Stdout)
	}
	var fileLogger *lj.Logger
	if filename != "" && strings.ToUpper(filename) != "STDOUT" {
		fileLogger = &lj.Logger{
			Filename:   filename,
			MaxSize:    100, // MB
			MaxBackups: 5,
			MaxAge:     365,   // days
			Compress:   false, // disabled by default
		}
		outputs = append(outputs, fileLogger)
	}
	logger := &logrus.Logger{
		Out:          io.MultiWriter(outputs...),
		Formatter:    &MyFormatter{PID: os.Getpid()},
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.Level(level),
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	return &KLLogger{
		withSTDOUT: withSTDOUT,
		level:      level,
		filename:   filename,
		logger:     logger,
		fileLogger: fileLogger,
	}
}

// KLLogger 日志对象
type KLLogger struct {
	withSTDOUT bool // 附加输出到stdout
	level      Level
	filename   string // 可以指定日志文件，或配置为"stdout"
	logger     *logrus.Logger
	fileLogger *lj.Logger // 写文件的对象
}

func (l *KLLogger) Close() {
	l.logger = nil
	if l.fileLogger != nil {
		l.fileLogger.Close()
		l.fileLogger = nil
	}
}

func (l *KLLogger) Debug(args ...interface{}) {
	if l.logger != nil {
		l.logger.Debug(args...)
	}
}

func (l *KLLogger) Info(args ...interface{}) {
	if l.logger != nil {
		l.logger.Info(args...)
	}
}

func (l *KLLogger) Warn(args ...interface{}) {
	if l.logger != nil {
		l.logger.Warn(args...)
	}
}

func (l *KLLogger) Error(args ...interface{}) {
	if l.logger != nil {
		l.logger.Error(args...)
	}
}

func (l *KLLogger) Panic(args ...interface{}) {
	if l.logger != nil {
		l.logger.Panic(args...)
	}
}

func (l *KLLogger) Fatal(args ...interface{}) {
	if l.logger != nil {
		l.logger.Fatal(args...)
	}
}

func (l *KLLogger) Debugf(template string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Debugf(template, args...)
	}
}

func (l *KLLogger) Infof(template string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Infof(template, args...)
	}
}

func (l *KLLogger) Warnf(template string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Warnf(template, args...)
	}
}

func (l *KLLogger) Errorf(template string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Errorf(template, args...)
	}
}

func (l *KLLogger) Fatalf(template string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Fatalf(template, args...)
	}
}

// io.Writer
func (l *KLLogger) Write(p []byte) (n int, err error) {
	l.Infof("%s", string(p))
	return len(p), nil
}
