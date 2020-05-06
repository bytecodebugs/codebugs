package log

import (
	"fmt"
	"github.com/fluent/fluent-logger-golang/fluent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger Logger
)

type Logger struct {
	*zap.Logger
}

func NewLogger(logger *zap.Logger) Logger {
	return Logger{logger}
}

type FileLogCfg struct {
	Enabled  bool          `json:"enabled"`
	Filename string        `json:"filename"`
	Level    zapcore.Level `json:"level"`
}
type SentryLogCfg struct {
	Enabled bool          `json:"enabled"`
	DSN     string        `json:"dsn"`
	Level   zapcore.Level `json:"level"`
}

type FluentLogCfg struct {
	Enabled bool          `json:"enabled"`
	Level   zapcore.Level `json:"level"`
	Config  fluent.Config `json:"config"`
}

type Config struct {
	FileLog   *FileLogCfg
	SentryLog *SentryLogCfg
	FluentLog *FluentLogCfg
}

func Init(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("log config can not be nil")
	}
	zapCfg := zap.NewProductionConfig()
	zapCfg.Sampling = nil
	zapLogger, err := zapCfg.Build()
	if err != nil {
		return fmt.Errorf("init logger error %q", err)
	}
	productionEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	var cores []zapcore.Core
	if cfg.FileLog.Enabled {
		fileLogger := NewFileLogger(cfg.FileLog.Filename)
		fileOut := zapcore.AddSync(fileLogger)
		fileCore := zapcore.NewCore(productionEncoder, fileOut, cfg.FileLog.Level)
		cores = append(cores, fileCore)
	}
	if cfg.SentryLog.Enabled {
		sentryCore, err := NewSentryCore(
			cfg.SentryLog.DSN,
			cfg.SentryLog.Level,
		)
		if err != nil {
			return fmt.Errorf("sentry logger init error %q", err)
		}
		cores = append(cores, sentryCore)
	}

	if cfg.FluentLog.Enabled {
		fluentCore, err := NewFluentCore(cfg.FluentLog.Config, cfg.FluentLog.Level)
		if err != nil {
			return fmt.Errorf("fluent logger init error %q", err)
		}
		cores = append(cores, fluentCore)
	}

	zapLogger = zapLogger.WithOptions(
		zap.WrapCore(func(zapcore.Core) zapcore.Core {
			return zapcore.NewTee(cores...)
		}),
	)
	logger = NewLogger(zapLogger)
	return nil
}

// Reopen log file when necessary.
// With creates a child logger and adds structured context to it.
func (logger Logger) With(fields ...zap.Field) Logger {
	return Logger{logger.Logger.With(fields...)}
}

// Named adds a new path segment to the logger's name.
func (logger Logger) Named(name string) Logger {
	return Logger{logger.Logger.Named(name)}
}

func Log() Logger {
	return logger
}
