package log

import "github.com/rs/zerolog"

const (
	LevelNone  = Level(zerolog.Disabled)
	LevelFatal = Level(zerolog.FatalLevel)
	LevelError = Level(zerolog.ErrorLevel)
	LevelInfo  = Level(zerolog.InfoLevel)
)

type Config struct {
	Level Level
}

type Level zerolog.Level
