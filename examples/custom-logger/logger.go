package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type FileLogger struct {
	logger *zap.Logger
}

type logLevel zapcore.Level

const (
	DebugLevel = logLevel(zapcore.DebugLevel)
	InfoLevel  = logLevel(zapcore.InfoLevel)
	WarnLevel  = logLevel(zapcore.WarnLevel)
	FatalLevel = logLevel(zapcore.FatalLevel)
)

// New creates a new FileLogger.
func New(logPath string, minLogLevel logLevel) (FileLogger, error) {
	file, err := os.Create(fmt.Sprintf("%s/botgo.log", logPath))
	if err != nil {
		return FileLogger{}, err
	}
	return FileLogger{
		logger: zap.New(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(file),
				zapcore.Level(minLogLevel),
			),
		),
	}, nil
}

// Debug logs a message at DebugLevel. The message includes any fields passed
func (f FileLogger) Debug(v ...interface{}) {
	f.logger.Debug(output(v...))
}

// Info logs a message at InfoLevel. The message includes any fields passed
func (f FileLogger) Info(v ...interface{}) {
	f.logger.Info(output(v...))
}

// Warn logs a message at WarnLevel. The message includes any fields passed
func (f FileLogger) Warn(v ...interface{}) {
	f.logger.Warn(output(v...))
}

// Error logs a message at ErrorLevel. The message includes any fields passed
func (f FileLogger) Error(v ...interface{}) {
	f.logger.Error(output(v...))
}

// Debugf logs a message at DebugLevel. The message includes any fields passed
func (f FileLogger) Debugf(format string, v ...interface{}) {
	f.logger.Debug(output(fmt.Sprintf(format, v...)))
}

// Infof logs a message at InfoLevel. The message includes any fields passed
func (f FileLogger) Infof(format string, v ...interface{}) {
	f.logger.Info(output(fmt.Sprintf(format, v...)))
}

// Warnf logs a message at WarnLevel. The message includes any fields passed
func (f FileLogger) Warnf(format string, v ...interface{}) {
	f.logger.Warn(output(fmt.Sprintf(format, v...)))
}

// Errorf logs a message at ErrorLevel. The message includes any fields passed
func (f FileLogger) Errorf(format string, v ...interface{}) {
	f.logger.Error(output(fmt.Sprintf(format, v...)))
}

// Sync flushes any buffered log entries.
func (f FileLogger) Sync() error {
	return f.logger.Sync()
}

func output(v ...interface{}) string {
	pc, file, line, _ := runtime.Caller(3)
	file = filepath.Base(file)
	funcName := strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), ".")

	logFormat := "%s %s:%d:%s " + fmt.Sprint(v...) + "\n"
	date := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf(logFormat, date, file, line, funcName)
}
