package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func init() {
	// Make the "time" field available for use
	zerolog.TimestampFieldName = "timestamp"
}

func New(cfg *Config) *Logger {
	out := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.DateTime
	})

	return &Logger{
		logger: zerolog.New(out).
			Level(zerolog.Level(cfg.Level)).With().
			Timestamp().
			Logger(),
	}
}

type Logger struct {
	logger zerolog.Logger
}

func (l *Logger) Fatal(msg string, err error) {
	l.log(LevelFatal, msg, err)
	os.Exit(1)
}

func (l *Logger) Error(msg string, err error) {
	l.log(LevelError, msg, err)
}

func (l *Logger) Info(msg string) {
	l.log(LevelInfo, msg, nil)
}

func (l *Logger) log(lvl Level, msg string, err error) {
	if l == nil {
		return
	}

	logEv := l.logger.WithLevel(zerolog.Level(lvl))
	if err != nil {
		logEv = logEv.Err(err)
	}
	logEv.Msg(msg)
}

func (l *Logger) WithFields(fields ...any) *Logger {
	if len(fields)%2 != 0 {
		panic("expected an even number of arguments")
	}

	if l == nil {
		return l
	}

	return &Logger{l.logger.With().Fields(fields).Logger()}
}
