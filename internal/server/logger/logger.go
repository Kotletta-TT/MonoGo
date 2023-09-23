package logger

import (
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var Logger *zap.SugaredLogger
var once sync.Once

func Init(config *config.Config) {
	once.Do(func() {
		var zapLevel zapcore.Level
		err := zapLevel.UnmarshalText([]byte(config.LogLevel))
		if err != nil {
			panic("invalid log level err: " + err.Error())
		}
		logFile, _, err := zap.Open(config.LogPath)
		if err != nil {
			panic("invalid log file err: " + err.Error())
		}
		stdoutCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.Lock(os.Stdout),
			zap.NewAtomicLevelAt(zapLevel))
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.Lock(logFile),
			zap.NewAtomicLevelAt(zapLevel))
		multiCore := zapcore.NewTee(stdoutCore, fileCore)
		Logger = zap.New(multiCore).Sugar()
	})
}
