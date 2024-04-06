package pkg

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	mode := viper.GetString("app.mode")
	switch mode {
	case "debug":
		zap.ReplaceGlobals(NewZapLogger(zapcore.DebugLevel))
	case "production":
		zap.ReplaceGlobals(NewZapLogger(zapcore.WarnLevel))
	default:
		zap.ReplaceGlobals(NewZapLogger(zapcore.InfoLevel))
	}
}

func NewZapLogger(level zapcore.Level) *zap.Logger {
	zapLogger, _ := zap.Config{
		Level:    zap.NewAtomicLevelAt(level),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			TimeKey:        "timestamp",
			LevelKey:       "level",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build(zap.AddCallerSkip(1))
	return zapLogger
}
