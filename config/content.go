package config

type RootConfig struct {
	Serve ServeCfg `json:"serve"`
	DB    DBCfg    `json:"db"`
	Log   LogCfg   `json:"log"`
}

type DBCfg struct {
}

type LogCfg struct {
	FileLog   FileLogCfg   `json:"file"`
	SentryLog SentryLogCfg `json:"sentry"`
	FluentLog FluentLogCfg `json:"fluent"`
}

type FileLogCfg struct {
	Enabled  bool   `json:"enabled"`
	Filename string `json:"filename"`
}
type SentryLogCfg struct {
	Enabled bool   `json:"enabled"`
	DSN     string `json:"dsn"`
}

type FluentLogCfg struct {
	Enabled bool   `json:"enabled"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Prefix  string `json:"prefix"`
}

type ServeCfg struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
