package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// InitLogger initializes the logger with the given service name
func InitLogger(serviceName string) *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder

	consoleEncoder := zapcore.NewJSONEncoder(config)

	// Create a core that writes to stdout
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
	)

	// Create the logger with some basic options
	log = zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Fields(
			zap.String("service", serviceName),
		),
	)

	return log
}

// GetLogger returns the logger instance
func GetLogger() *zap.Logger {
	if log == nil {
		return InitLogger("unknown")
	}
	return log
}

// Info logs a message at InfoLevel
func Info(message string, fields ...zap.Field) {
	GetLogger().Info(message, fields...)
}

// Debug logs a message at DebugLevel
func Debug(message string, fields ...zap.Field) {
	GetLogger().Debug(message, fields...)
}

// Warn logs a message at WarnLevel
func Warn(message string, fields ...zap.Field) {
	GetLogger().Warn(message, fields...)
}

// Error logs a message at ErrorLevel
func Error(message string, fields ...zap.Field) {
	GetLogger().Error(message, fields...)
}

// Fatal logs a message at FatalLevel and then calls os.Exit(1)
func Fatal(message string, fields ...zap.Field) {
	GetLogger().Fatal(message, fields...)
}

// With creates a child logger with the given fields
func With(fields ...zap.Field) *zap.Logger {
	return GetLogger().With(fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return GetLogger().Sync()
}
