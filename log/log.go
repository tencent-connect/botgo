// Package log 是 SDK 的 logger 接口定义与内置的 logger。
package log

// DefaultLogger 默认logger
var DefaultLogger = Logger(new(consoleLogger))

// Debug log.Debug
func Debug(v ...interface{}) {
	DefaultLogger.Debug(v...)
}

// Info log.Info
func Info(v ...interface{}) {
	DefaultLogger.Info(v...)
}

// Warn log.Warn
func Warn(v ...interface{}) {
	DefaultLogger.Warn(v...)
}

// Error log.Error
func Error(v ...interface{}) {
	DefaultLogger.Error(v...)
}

// Debugf log.Debugf
func Debugf(format string, v ...interface{}) {
	DefaultLogger.Debugf(format, v...)
}

// Infof log.Infof
func Infof(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

// Warnf log.Warnf
func Warnf(format string, v ...interface{}) {
	DefaultLogger.Warnf(format, v...)
}

// Errorf log.Errorf
func Errorf(format string, v ...interface{}) {
	DefaultLogger.Errorf(format, v...)
}

// Sync logger Sync calls to flush buffer
func Sync() {
	_ = DefaultLogger.Sync()
}
