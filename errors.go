package web

// Logger defines all the logging methods to be implemented
type Logger interface {
	Debug(data ...interface{})
	Info(data ...interface{})
	Warn(data ...interface{})
	Error(data ...interface{})
	Fatal(data ...interface{})
}
