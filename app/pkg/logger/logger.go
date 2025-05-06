package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type LoggerInterface interface {
	GetLogger() *zerolog.Logger
	WithFields(fields map[string]interface{}) LoggerInterface
	WithContext(requestID, method string) LoggerInterface
	AddFields(key, val string) LoggerInterface
	Info(msg string, fields map[string]interface{})
	Error(err error, msg string, fields map[string]interface{})
	Fatal(err error, msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{}, err error)
	Debug(msg string, fields map[string]interface{})
	WrapError(message string, err error) error
	WrapDetailError(message string, err error, messDetail interface{}) error
}

func ConfigureLogger(env string) *zerolog.Logger {
	var logger zerolog.Logger
	timeFormat := "15:04:05"

	switch env {
	case "envLocal":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: timeFormat})

	case "envDev":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	case "envProd":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true})

	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	logger = logger.With().Timestamp().Logger()

	return &logger
}

type Logger struct {
	Logger *zerolog.Logger
}

func NewLogger(logger *zerolog.Logger) *Logger {
	return &Logger{Logger: logger}
}

func (l *Logger) GetLogger() *zerolog.Logger {
	return l.Logger
}

// WithFields создает новый логгер с предварительно установленными полями.
func (l *Logger) WithFields(fields map[string]interface{}) LoggerInterface {
	context := l.Logger.With().Fields(fields).Logger()
	l.Logger = &context

	return l
}

func (l *Logger) WithContext(requestID, method string) LoggerInterface {
	return l.WithFields(map[string]interface{}{
		"request_id": requestID,
		"method":     method,
	})
}

func (l *Logger) AddFields(key string, val string) LoggerInterface {
	context := l.Logger.With().Str(key, val).Logger()
	return &Logger{Logger: &context}
}

func (l *Logger) Info(msg string, fields map[string]interface{}) {
	l.Logger.Info().Fields(fields).Msg(msg)
}

func (l *Logger) Error(err error, msg string, fields map[string]interface{}) {
	l.Logger.Error().Err(err).Fields(fields).Msg(msg)
}

func (l *Logger) Fatal(err error, msg string, fields map[string]interface{}) {
	l.Logger.Fatal().Err(err).Fields(fields).Msg(msg)
}

func (l *Logger) Warn(msg string, fields map[string]interface{}, err error) {
	l.Logger.Warn().Err(err).Fields(fields).Msg(msg)
}

func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	l.Logger.Debug().Fields(fields).Msg(msg)
}

func (l *Logger) WrapError(message string, err error) error {
	l.Logger.Error().Err(err).Msg(message)

	return fmt.Errorf("%s", message)
}

func (l *Logger) WrapDetailError(message string, err error, messDetail interface{}) error {
	l.Logger.Error().Err(err).Interface("messDetail", messDetail).Msg(message)

	return fmt.Errorf("%s", message)
}
