package logger_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errPlain = errors.New("connection reset")

type loggedError struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"` //nolint:tagliatelle // log fields are snake_case
}

func marshal(t *testing.T, we *logger.WideEvent) loggedError {
	t.Helper()

	var buf bytes.Buffer
	l := logger.NewProductionLogger(&buf)
	l.Info().EmbedObject(we).Send()

	var out loggedError
	require.NoError(t, json.Unmarshal(buf.Bytes(), &out))

	return out
}

func Test_WideEvent_ErrorCode(t *testing.T) {
	t.Run("cron job with an ApiError logs its code and the full message", func(t *testing.T) {
		we := &logger.WideEvent{}
		we.FinishCronJob(errs.ProductNotFound("p1"), 0)

		out := marshal(t, we)
		assert.Equal(t, errs.ErrProductNotFound.Code, out.ErrorCode)
		assert.Contains(t, out.Error, "p1")
	})

	t.Run("cron job with a plain error logs INTERNAL_ERROR", func(t *testing.T) {
		we := &logger.WideEvent{}
		we.FinishCronJob(errPlain, 0)

		out := marshal(t, we)
		assert.Equal(t, errs.ErrInternal.Code, out.ErrorCode)
		assert.Equal(t, errPlain.Error(), out.Error)
	})

	t.Run("request logs the given code even without an error to log", func(t *testing.T) {
		we := &logger.WideEvent{}
		we.FinishRequest(400, "VALIDATION_ERROR", nil, 0)

		out := marshal(t, we)
		assert.Equal(t, "VALIDATION_ERROR", out.ErrorCode)
		assert.Empty(t, out.Error)
	})
}
