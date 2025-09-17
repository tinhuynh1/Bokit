package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.Logger

func Init() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	var err error
	L, err = config.Build()
	if err != nil {
		panic(err)
	}
}
