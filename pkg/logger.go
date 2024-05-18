package pkg

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
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
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename: "github_observer.log",
		MaxSize:  10,
		MaxAge:   3,
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
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
		}),
		w,
		level,
	)

	zapLogger := zap.New(core, zap.AddCallerSkip(1))
	return zapLogger
}
