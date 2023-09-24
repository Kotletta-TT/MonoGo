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
		Logger = zap.New(multiCore).Sugar()
	})
}
