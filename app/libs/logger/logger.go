package logger

import (
	"context"
	"os"
	"rabi-food-core/config"

	"github.com/rs/zerolog"
)

var (
	base      zerolog.Logger
	LoggerKey = "loggerKey"
)

func Init() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if config.Env == "dev" {
		base = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	} else {
		base = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}

// Get retorna o logger da requisição atual
func Get(c context.Context) *zerolog.Logger {
	v := c.Value(LoggerKey)
	if l, ok := v.(*zerolog.Logger); ok {
		return l
	}
	return L()
}

func L() *zerolog.Logger {
	return &base
}
