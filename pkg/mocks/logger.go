package mocks

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func MockedLogger() *observer.ObservedLogs {
	core, logs := observer.New(zap.DebugLevel)
	zapLogger := zap.New(core)
	zap.ReplaceGlobals(zapLogger)
	return logs
}
