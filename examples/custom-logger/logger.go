package main

import (
	"fmt"
	"os"
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
			)),
	}, nil
}

func (f FileLogger) Debug(v ...interface{}) {
	f.logger.Debug(output(v...))
}

func (f FileLogger) Info(v ...interface{}) {
	f.logger.Info(output(v...))
}

func (f FileLogger) Warn(v ...interface{}) {
	f.logger.Warn(output(v...))
}

func (f FileLogger) Error(v ...interface{}) {
	f.logger.Error(output(v...))
}

func (f FileLogger) Debugf(format string, v ...interface{}) {
	f.logger.Debug(output(fmt.Sprintf(format, v...)))
}

func (f FileLogger) Infof(format string, v ...interface{}) {
	f.logger.Info(output(fmt.Sprintf(format, v...)))
}

func (f FileLogger) Warnf(format string, v ...interface{}) {
	f.logger.Warn(output(fmt.Sprintf(format, v...)))
}

func (f FileLogger) Errorf(format string, v ...interface{}) {
	f.logger.Error(output(fmt.Sprintf(format, v...)))
}

func (f FileLogger) Sync() error {
	return f.logger.Sync()
}

func output(v ...interface{}) string {
	_, file, line, _ := runtime.Caller(3)
	files := strings.Split(file, "/")
	file = files[len(files)-1]

	logFormat := "%s %s:%d " + fmt.Sprint(v...) + "\n"
	date := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf(logFormat, date, file, line)
}
