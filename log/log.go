package log

var gLogger *KLLogger
var gDefaultLogger *KLLogger

func init() {
	// 未初始化时，默认的日志输出
	gDefaultLogger = NewKLLogger(InfoLevel, "", true)
	gLogger = gDefaultLogger
}

// Init log service
func Init(debugMode bool, filename string) {
	if gLogger != nil && gLogger != gDefaultLogger {
		gLogger.Fatal("logger already initialized")
	}
	if debugMode {
		gLogger = NewKLLogger(DebugLevel, filename, true)
	} else {
		gLogger = NewKLLogger(InfoLevel, filename, false)
	}
}

// GetKLLogger .
func GetKLLogger() *KLLogger {
	return gLogger
}

// Close .
func Close() {
	if gLogger != nil {
		gLogger.Close()
	}
}

func Debug(args ...interface{}) {
	gLogger.Debug(args...)
}

func Info(args ...interface{}) {
	gLogger.Info(args...)
}

func Warn(args ...interface{}) {
	gLogger.Warn(args...)
}

func Error(args ...interface{}) {
	gLogger.Error(args...)
}

func Panic(args ...interface{}) {
	gLogger.Panic(args...)
}

func Fatal(args ...interface{}) {
	gLogger.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	gLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	gLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	gLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	gLogger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	gLogger.Fatalf(template, args...)
}
