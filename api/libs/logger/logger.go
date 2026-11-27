package logger

import (
	"context"
	"io"
	"os"
	"strings"
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

	json, formatValid := writesJSON()
	base = NewWithOutput(os.Stdout, json)
	if !formatValid {
		// Like LOG_LEVEL: a typo keeps the environment's default and says so.
		base.Warn().Str("log_format", config.LogFormat).Msg("invalid LOG_FORMAT, using the environment default")
	}

	level, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil || level == zerolog.NoLevel {
		// A typo in LOG_LEVEL must not take the API down: keep the default and say so.
		level = zerolog.InfoLevel
		base.Warn().Str("log_level", config.LogLevel).Msg("invalid LOG_LEVEL, using info")
	}

	zerolog.SetGlobalLevel(level)
}

// writesJSON reports whether the logs are written as JSON lines, the format the log pipeline parses, and
// whether LOG_FORMAT was a known value.
func writesJSON() (bool, bool) {
	switch strings.ToLower(config.LogFormat) {
	case "json":
		return true, true
	case "text":
		return false, true
	case "":
		return config.Env.IsProduction(), true
	default:
		return config.Env.IsProduction(), false
	}
}

// NewWithOutput builds the application logger writing to out: JSON lines, the format the log pipeline parses,
// or readable text. The app uses it once, for os.Stdout; tests use it to capture the production log.
func NewWithOutput(out io.Writer, json bool) zerolog.Logger {
	if json {
		return newLogger(out)
	}

	return newLogger(zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: time.RFC3339,
		NoColor:    true,
	})
}

// SetOutput sends the logs to out, in the production format. Tests use it to read the log the API
// emits; call it before the API serves requests.
func SetOutput(out io.Writer) {
	base = NewWithOutput(out, true)
}

// newLogger holds what every environment shares: only the writer differs between them.
func newLogger(writer io.Writer) zerolog.Logger {
	return zerolog.New(writer).With().Timestamp().Str("version", config.ApiVersion).Logger()
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
