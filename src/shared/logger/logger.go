package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func InitLogger(serviceName string) *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder

	consoleEncoder := zapcore.NewJSONEncoder(config)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
	)

	log = zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Fields(
			zap.String("service", serviceName),
		),
	)

	return log
}

func GetLogger() *zap.Logger {
	if log == nil {
		return InitLogger("unknown")
	}
	return log
}

func Info(message string, fields ...zap.Field) {
	GetLogger().Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	GetLogger().Debug(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	GetLogger().Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	GetLogger().Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	GetLogger().Fatal(message, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return GetLogger().With(fields...)
}

func Sync() error {
	return GetLogger().Sync()
}
