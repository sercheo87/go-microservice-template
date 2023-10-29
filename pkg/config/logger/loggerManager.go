package logger

import (
	"context"
	"fmt"
	"go-microservice-template/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

type RequestIDKey struct{}

type ctxKey struct{}

var LoggerInstance *zap.Logger

var Ctx context.Context

var core zapcore.Core

func ApplyLoggerConfiguration(logLevel string) (*zap.Logger, error) {
	var level zap.AtomicLevel

	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	}

	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	fileEncoder := zapcore.NewJSONEncoder(zapEncoderConfig)
	workspaceFolder, _ := filepath.Abs(config.Configuration.WorkspaceFolder)

	logFile := fmt.Sprintf("%s/logs/%s-%s.log", workspaceFolder, config.Configuration.MicroserviceName, config.Configuration.Environment)

	logRotation := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // Max size in megabytes before log rotation
		MaxBackups: 10, // Max number of old log files to keep
		MaxAge:     30, // Max number of days to retain old log files
		LocalTime:  true,
		Compress:   true, // Compress old log files
	}

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	defaultLogLevel := level
	core = zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logRotation), defaultLogLevel).With(
			[]zap.Field{
				zap.String("app", config.Configuration.MicroserviceName),
				zap.String("env", config.Configuration.Environment),
			}),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), defaultLogLevel).With(
			[]zap.Field{
				zap.String("app", config.Configuration.MicroserviceName),
				zap.String("env", config.Configuration.Environment),
			}),
	)

	LoggerInstance = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return LoggerInstance, nil
}

func FromCtx() *zap.Logger {
	if l, ok := Ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	} else if l := LoggerInstance; l != nil {
		return l
	}

	return zap.NewNop()
}

func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	Ctx = ctx
	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		if lp == l {
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, l)
}

func Reset() {
	Ctx = nil
	clone := core.With([]zap.Field{})
	LoggerInstance = zap.New(clone)
}
