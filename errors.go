package web

import (
	"errors"
	"io"
	"log"
)

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

var (
	// ErrInvalidPort is the error returned when the port number provided in the config file is invalid
	ErrInvalidPort = errors.New("Port number not provided or is invalid (should be between 0 - 65535)")
	lh             *logHandler
	// LOGHANDLER is a global variable which web uses to log messages
	LOGHANDLER Logger
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

// logHandler has all the log writer handlers
type logHandler struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
	fatal *log.Logger
}

// Debug prints log of severity 5
func (lh *logHandler) Debug(data ...interface{}) {
	if lh.debug == nil {
		return
	}
	lh.debug.Println(data...)
}

// Info prints logs of severity 4
func (lh *logHandler) Info(data ...interface{}) {
	if lh.info == nil {
		return
	}
	lh.info.Println(data...)
}

// Warn prints log of severity 3
func (lh *logHandler) Warn(data ...interface{}) {
	if lh.warn == nil {
		return
	}
	lh.warn.Println(data...)
}

// Error prints log of severity 2
func (lh *logHandler) Error(data ...interface{}) {
	if lh.err == nil {
		return
	}
	lh.err.Println(data...)
}

// Fatal prints log of severity 1
func (lh *logHandler) Fatal(data ...interface{}) {
	if lh.fatal == nil {
		return
	}
	lh.fatal.Fatalln(data...)
}

func loggerWithCfg(stdout io.Writer, stderr io.Writer, cfgs ...logCfg) *logHandler {
	lh = &logHandler{
		debug: log.New(stdout, "Debug ", log.LstdFlags),
		info:  log.New(stdout, "Info ", log.LstdFlags),
		warn:  log.New(stderr, "Warning ", log.LstdFlags),
		err:   log.New(stderr, "Error ", log.LstdFlags),
		fatal: log.New(stderr, "Fatal ", log.LstdFlags|log.Llongfile),
	}

	for _, c := range cfgs {
		switch c {
		case LogCfgDisableDebug:
			{
				lh.debug = nil
			}
		case LogCfgDisableInfo:
			{
				lh.info = nil
			}
		case LogCfgDisableWarn:
			{
				lh.warn = nil
			}
		case LogCfgDisableError:
			{
				lh.err = nil
			}
		case LogCfgDisableFatal:
			{
				lh.fatal = nil
			}
		}
	}
	return lh
}
