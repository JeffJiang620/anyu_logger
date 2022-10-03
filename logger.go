package loggers

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logLevelMap = map[string]zapcore.Level{
	"info":  zap.InfoLevel,
	"debug": zap.DebugLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
}

func GetLogLevel(level string) zapcore.Level {
	if logLevel, ok := logLevelMap[level]; ok {
		return logLevel
	} else {
		return zap.InfoLevel
	}
}

var Logger = NewCommonLogger(DefaultLoggerConfig).Options(zap.AddCaller())

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

var defaultEncodeConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "name",
	CallerKey:      "caller",
	FunctionKey:    "",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     timeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

var DefaultLoggerConfig = &Config{
	LogLevel:     "info",
	Writer:       os.Stdout,
	EncodeConfig: defaultEncodeConfig,
}

type Config struct {
	LogLevel     string
	Writer       io.Writer
	EncodeConfig zapcore.EncoderConfig
}

type CommonLogger struct {
	*zap.Logger
	config *Config
}

func NewCommonLogger(config *Config) *CommonLogger {
	syncer := zapcore.AddSync(config.Writer)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(config.EncodeConfig), syncer, GetLogLevel(config.LogLevel))
	return &CommonLogger{
		Logger: zap.New(core),
		config: config,
	}
}

func (logger *CommonLogger) clone() *CommonLogger {
	c := *logger
	return &c
}

func (logger *CommonLogger) LogLevel(logLevel string) *CommonLogger {
	newLogger := logger.clone()
	newLogger.Logger = newLogger.WithOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		syncer := zapcore.AddSync(logger.config.Writer)
		core := zapcore.NewCore(zapcore.NewJSONEncoder(logger.config.EncodeConfig), syncer, GetLogLevel(logLevel))
		return core
	}))
	return newLogger
}

func (logger *CommonLogger) EncodeConfig(encodeConfig zapcore.EncoderConfig) *CommonLogger {
	newLogger := logger.clone()
	newLogger.Logger = newLogger.WithOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		syncer := zapcore.AddSync(logger.config.Writer)
		core := zapcore.NewCore(zapcore.NewJSONEncoder(encodeConfig), syncer, GetLogLevel(logger.config.LogLevel))
		return core
	}))
	return newLogger
}

func (logger *CommonLogger) Writer(writer io.Writer) *CommonLogger {
	newLogger := logger.clone()
	newLogger.Logger = newLogger.WithOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		syncer := zapcore.AddSync(writer)
		core := zapcore.NewCore(zapcore.NewJSONEncoder(logger.config.EncodeConfig), syncer, GetLogLevel(logger.config.LogLevel))
		return core
	}))
	return newLogger
}

func (logger *CommonLogger) Options(options ...zap.Option) *CommonLogger {
	newLogger := logger.clone()
	newLogger.Logger = newLogger.WithOptions(options...)
	return newLogger
}

func (logger *CommonLogger) Name(name string) *CommonLogger {
	newLogger := logger.clone()
	newLogger.Logger = newLogger.Named(name)
	return newLogger
}
