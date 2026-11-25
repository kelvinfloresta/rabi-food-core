package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/Rabi-IT/rabi-food-core/config"

	"github.com/rs/zerolog"
)

var base zerolog.Logger

const (
	loggerKey = "loggerKey"
)

func init() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	base = newBase(os.Stdout, config.Env == "production")
}

// newBase builds the application logger. Production writes JSON; other environments write readable text.
func newBase(out io.Writer, production bool) zerolog.Logger {
	if !production {
		out = zerolog.ConsoleWriter{
			Out:        out,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		}
	}

	return zerolog.New(out).With().Timestamp().Str("version", config.ApiVersion).Logger()
}

// Get retrieves the logger from the given context.
// If no logger is found, returns a pointer to a fresh copy of the base logger
// so that callers can safely call UpdateContext without mutating global state.
func Get(c context.Context) *zerolog.Logger {
	v := c.Value(loggerKey)
	if l, ok := v.(*zerolog.Logger); ok {
		return l
	}

	l := base.With().Logger()

	return &l
}

func New() zerolog.Context {
	return base.With()
}

// WithContext stores the logger in the context.
func WithContext(c context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(c, loggerKey, logger)
}
