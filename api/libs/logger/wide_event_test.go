package logger_test

import (
	"errors"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/fixtures"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/stretchr/testify/assert"
)

var errPlain = errors.New("connection reset")

func Test_WideEvent_ErrorCode(t *testing.T) {
	t.Run("cron job with an ApiError logs its code and the full message", func(t *testing.T) {
		we := &logger.WideEvent{}
		we.FinishCronJob(errs.ProductNotFound("p1"), 0)

		fields := fixtures.Logger.WideEvent(t, we)
		assert.Equal(t, errs.ErrProductNotFound.Code, fields["error_code"])
		assert.Contains(t, fields["error"], "p1")
	})

	t.Run("cron job with a plain error logs INTERNAL_ERROR", func(t *testing.T) {
		we := &logger.WideEvent{}
		we.FinishCronJob(errPlain, 0)

		fields := fixtures.Logger.WideEvent(t, we)
		assert.Equal(t, errs.ErrInternal.Code, fields["error_code"])
		assert.Equal(t, errPlain.Error(), fields["error"])
	})

	t.Run("request logs request_id, user_id and selected_tenant_id to find a reported failure", func(t *testing.T) {
		we := &logger.WideEvent{RequestID: "req-1"}
		we.SetUserID("user-1")
		we.SetSelectedTenantID("tenant-1")
		we.FinishRequest(500, errs.ErrInternal.Code, errPlain, 0)

		fields := fixtures.Logger.WideEvent(t, we)
		assert.Equal(t, "req-1", fields["request_id"])
		assert.Equal(t, "user-1", fields["user_id"])
		assert.Equal(t, "tenant-1", fields["selected_tenant_id"])
	})

	t.Run("request logs the given code even without an error to log", func(t *testing.T) {
		we := &logger.WideEvent{}
		we.FinishRequest(400, "VALIDATION_ERROR", nil, 0)

		fields := fixtures.Logger.WideEvent(t, we)
		assert.Equal(t, "VALIDATION_ERROR", fields["error_code"])
		assert.NotContains(t, fields, "error")
	})
}
