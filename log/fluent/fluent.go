package fluent

import (
	"os"

	"github.com/fluent/fluent-logger-golang/fluent"
	"go.uber.org/zap/zapcore"
)

const (
	FATAL = "F"
	ERROR = "E"
	WARN  = "W"
	INFO  = "I"
	DEBUG = "D"
)

type Config struct {
	Fluent fluent.Config `json:"fluent"`
}

type Packet struct {
	Timestamp string                 `msg:"timestamp"`
	Level     string                 `msg:"convertLevel"`
	Host      string                 `msg:"host"`
	Message   string                 `msg:"message"`
	Extra     map[string]interface{} `msg:"extra"`
	Caller    zapcore.EntryCaller    `msg:"caller"`
}

func convertLevel(lvl zapcore.Level) string {
	switch lvl {
	case zapcore.DebugLevel:
		return DEBUG
	case zapcore.InfoLevel:
		return INFO
	case zapcore.WarnLevel:
		return WARN
	case zapcore.ErrorLevel:
		return ERROR
	case zapcore.DPanicLevel:
		return FATAL
	case zapcore.PanicLevel:
		return FATAL
	case zapcore.FatalLevel:
		return FATAL
	default:
		return FATAL
	}
}

var hostName string

func init() {
	hostName, _ = os.Hostname()
}

func New(cfg Config, enabler zapcore.LevelEnabler) (zapcore.Core, error) {
	client, err := fluent.New(cfg.Fluent)
	return &Logger{
		LevelEnabler: enabler,
		client:       client,
		fields:       make(map[string]interface{}),
	}, err
}

type Logger struct {
	zapcore.LevelEnabler
	client *fluent.Fluent
	fields map[string]interface{}
}

func (core *Logger) With(fields []zapcore.Field) zapcore.Core {
	return core.with(fields)
}

func (core *Logger) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if core.Enabled(ent.Level) {
		return ce.AddCore(ent, core)
	}
	return ce
}

func (core *Logger) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	clone := core.with(fields)
	packet := &Packet{
		Timestamp: ent.Time.Format("2006-01-02T15:04:05"),
		Host:      hostName,
		Message:   ent.Message,
		Level:     convertLevel(ent.Level),
		Extra:     clone.fields,
		Caller:    ent.Caller,
	}
	err := core.Post(packet)
	return err
}
func (core *Logger) Post(packet *Packet) error {
	err := core.client.Post("server.log", *packet)
	return err
}
func (core *Logger) Sync() error {
	return nil
}

func (core *Logger) with(fields []zapcore.Field) *Logger {
	m := make(map[string]interface{}, len(core.fields))
	for k, v := range core.fields {
		m[k] = v
	}

	// Add fields to an in-memory encoder.
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}

	// Merge the two maps.
	for k, v := range enc.Fields {
		m[k] = v
	}

	return &Logger{
		LevelEnabler: core.LevelEnabler,
		client:       core.client,
		fields:       m,
	}
}
