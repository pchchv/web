package web

const (
	// LogCfgDisableDebug is used to disable debug logs
	LogCfgDisableDebug = logCfg("disable-debug")
	// LogCfgDisableInfo is used to disable info logs
	LogCfgDisableInfo = logCfg("disable-info")
	// LogCfgDisableWarn is used to disable warning logs
	LogCfgDisableWarn = logCfg("disable-warn")
	// LogCfgDisableError is used to disable error logs
	LogCfgDisableError = logCfg("disable-err")
	// LogCfgDisableFatal is used to disable fatal logs
	LogCfgDisableFatal = logCfg("disable-fatal")
)

type logCfg string

// Logger defines all the logging methods to be implemented
type Logger interface {
	Debug(data ...interface{})
	Info(data ...interface{})
	Warn(data ...interface{})
	Error(data ...interface{})
	Fatal(data ...interface{})
}
