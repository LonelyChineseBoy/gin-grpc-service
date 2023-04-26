package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetEncoder() zapcore.Encoder {
	zap.NewProductionConfig()
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "time"
	config.NameKey = "name"
	config.CallerKey = "caller"
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeDuration = zapcore.SecondsDurationEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder
	config.LineEnding = zapcore.DefaultLineEnding
	return zapcore.NewJSONEncoder(config)
}
