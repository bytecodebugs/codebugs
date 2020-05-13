package log

import (
	"codebugs/log/file"
	"codebugs/log/fluent"
	"codebugs/log/sentry"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

type FileCfg struct {
	Enabled  bool          `json:"enabled"`
	Level    zapcore.Level `json:"level"`
	Filename string        `json:"filename"`
}

type SentryCfg struct {
	Enabled bool          `json:"enabled"`
	Level   zapcore.Level `json:"level"`
	DSN     string        `json:"dsn"`
}

type FluentCfg struct {
	Enabled bool          `json:"enabled"`
	Level   zapcore.Level `json:"level"`
	Config  fluent.Config `json:"config"`
}

type Config struct {
	FileLog   *FileCfg
	SentryLog *SentryCfg
	FluentLog *FluentCfg
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
		fileLogger := file.New(cfg.FileLog.Filename)
		fileOut := zapcore.AddSync(fileLogger)
		fileCore := zapcore.NewCore(productionEncoder, fileOut, cfg.FileLog.Level)
		cores = append(cores, fileCore)
	}
	if cfg.SentryLog.Enabled {
		sentryCore, err := sentry.New(
			cfg.SentryLog.DSN,
			cfg.SentryLog.Level,
		)
		if err != nil {
			return fmt.Errorf("sentry logger init error %q", err)
		}
		cores = append(cores, sentryCore)
	}

	if cfg.FluentLog.Enabled {
		fluentCore, err := fluent.New(cfg.FluentLog.Config, cfg.FluentLog.Level)
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
	logger = zapLogger
	return nil
}

func Log() *zap.Logger {
	return logger
}
