// Package logger provides a logger for the application.
package logger

import (
	"os"
	"sync"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger
var once sync.Once

// Init initializes the application with the given configuration.
//
// It takes a pointer to a config.Config struct as a parameter.
// It does not return anything.
func Init(config *config.Config) {
	once.Do(func() {
		cores := make([]zapcore.Core, 0, 2)
		var zapLevel zapcore.Level
		err := zapLevel.UnmarshalText([]byte(config.LogLevel))
		if err != nil {
			panic("invalid log level err: " + err.Error())
		}
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.Lock(os.Stdout),
			zap.NewAtomicLevelAt(zapLevel)))
		if config.LogFile {
			logFile, _, err := zap.Open(config.LogPath)
			if err != nil {
				panic("invalid log file err: " + err.Error())
			}
			cores = append(cores, zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.Lock(logFile),
				zap.NewAtomicLevelAt(zapLevel)))
		}
		multiCore := zapcore.NewTee(cores...)
		logger = zap.New(multiCore).Sugar()
	})
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func DPanic(args ...interface{}) {
	logger.DPanic(args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func DPanicf(template string, args ...interface{}) {
	logger.DPanicf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	logger.DPanicw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	logger.Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	logger.Fatalw(msg, keysAndValues...)
}

func Debugln(args ...interface{}) {
	logger.Debugln(args...)
}

func Infoln(args ...interface{}) {
	logger.Infoln(args...)
}

func Warnln(args ...interface{}) {
	logger.Warnln(args...)
}

func Errorln(args ...interface{}) {
	logger.Errorln(args...)
}

func DPanicln(args ...interface{}) {
	logger.DPanicln(args...)
}

func Panicln(args ...interface{}) {
	logger.Panicln(args...)
}

func Fatalln(args ...interface{}) {
	logger.Fatalln(args...)
}

func Sync() error {
	return logger.Sync()
}
