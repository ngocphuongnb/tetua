package logger

type Context map[string]interface{}

type Logger interface {
	Info(...interface{})
	Error(...interface{})
	Debug(...interface{})
	Fatal(...interface{})
	Warn(...interface{})
	Panic(...interface{})
	DPanic(...interface{})
	WithContext(context Context) Logger
}

var loggerInstance Logger

func New(logger Logger) {
	loggerInstance = logger
}

func Get() Logger {
	return loggerInstance
}

func Debug(params ...interface{}) {
	loggerInstance.Debug(params...)
}

func Info(params ...interface{}) {
	loggerInstance.Info(params...)
}

func Warn(params ...interface{}) {
	loggerInstance.Warn(params...)
}

func Error(params ...interface{}) {
	loggerInstance.Error(params...)
}

func DPanic(params ...interface{}) {
	loggerInstance.DPanic(params...)
}

func Panic(params ...interface{}) {
	loggerInstance.Panic(params...)
}

func Fatal(params ...interface{}) {
	loggerInstance.Fatal(params...)
}
