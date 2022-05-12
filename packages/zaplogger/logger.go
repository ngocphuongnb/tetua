package zapLogger

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/ngocphuongnb/tetua/app/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	*zap.Logger
	logger.Context
}

type Config struct {
	Development bool `json:"development"`
	LogFile     string
}

var loggerInstance *ZapLogger

func New(config Config) *ZapLogger {
	if config.LogFile != "" {
		if err := os.MkdirAll(path.Dir(config.LogFile), 0755); err != nil {
			log.Fatal(err)
		}
		logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		logFile.Close()
	}
	zapConfig := zap.NewProductionEncoderConfig()
	// zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)
	fileEncoder := zapcore.NewJSONEncoder(zapConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(zapConfig)
	logFile, _ := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	zapLogger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	defer zapLogger.Sync()
	return &ZapLogger{zapLogger, logger.Context{}}
}

func Get(context logger.Context) logger.Logger {
	return &ZapLogger{loggerInstance.Logger, context}
}

func WithContext(contexts ...logger.Context) logger.Logger {
	if len(contexts) > 0 {
		return &ZapLogger{loggerInstance.Logger, contexts[0]}
	}

	return loggerInstance
}

func (l *ZapLogger) WithContext(context logger.Context) logger.Logger {
	return &ZapLogger{l.Logger, context}
}

func (l *ZapLogger) Debug(params ...interface{}) {
	msg, contexts := l.getLogContext(params...)
	l.Logger.Debug(msg, getZapFields(contexts...)...)
}

func (l *ZapLogger) Info(params ...interface{}) {
	msg, contexts := l.getLogContext(params...)
	l.Logger.Info(msg, getZapFields(contexts...)...)
}

func (l *ZapLogger) Warn(params ...interface{}) {
	msg, contexts := l.getLogContext(params...)
	l.Logger.Warn(msg, getZapFields(contexts...)...)
}

func (l *ZapLogger) Error(params ...interface{}) {
	msg, contexts := l.getLogContext(params...)
	l.Logger.Error(msg, getZapFields(contexts...)...)
}

func (l *ZapLogger) DPanic(params ...interface{}) {
	msg, contexts := l.getLogContext(params...)
	l.Logger.DPanic(msg, getZapFields(contexts...)...)
}

func (l *ZapLogger) Panic(params ...interface{}) {
	msg, contexts := l.getLogContext(params...)
	l.Logger.Panic(msg, getZapFields(contexts...)...)
}

func (l *ZapLogger) Fatal(params ...interface{}) {
	msg, contexts := l.getLogContext(params...)
	l.Logger.Fatal(msg, getZapFields(contexts...)...)
}

func getZapFields(contexts ...logger.Context) []zapcore.Field {
	var contextFields []zapcore.Field
	for _, context := range contexts {
		for key, val := range context {
			// keyIndex := fmt.Sprintf("%d_%s", contexIndex, key)
			keyIndex := key
			contextFields = append(contextFields, zap.Any(keyIndex, val))
		}
	}
	return contextFields
}
