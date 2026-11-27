package logger_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/fixtures"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

var errPlain = errors.New("connection reset")

func Test_WideEvent_ErrorCode(t *testing.T) {
	t.Run("cron job with an ApiError logs its code and the full message", func(t *testing.T) {
		we := &logger.WideEvent{}
		we.FinishCronJob(errs.ProductNotFound("p1"), 0)

		fields := fixtures.Logger.WideEvent(t, we)
		assert.Equal(t, errs.ErrProductNotFound.Code, fields["error.type"])
		assert.Contains(t, fields["exception.message"], "p1")
	})

	t.Run("cron job with a plain error logs INTERNAL_ERROR", func(t *testing.T) {
		we := &logger.WideEvent{}
		we.FinishCronJob(errPlain, 0)

		fields := fixtures.Logger.WideEvent(t, we)
		assert.Equal(t, errs.ErrInternal.Code, fields["error.type"])
		assert.Equal(t, errPlain.Error(), fields["exception.message"])
	})

	t.Run("request logs request_id, auth_user_id and selected_tenant_id to find a reported failure", func(t *testing.T) {
		we := &logger.WideEvent{RequestID: "req-1"}
		we.SetAuthUserID("user-1")
		we.SetSelectedTenantID("tenant-1")
		we.FinishRequest(500, errs.ErrInternal.Code, errPlain, 0)

		fields := fixtures.Logger.WideEvent(t, we)
		assert.Equal(t, "req-1", fields["request_id"])
		assert.Equal(t, "user-1", fields["auth_user_id"])
		assert.Equal(t, "tenant-1", fields["selected_tenant_id"])
	})

	t.Run("request logs the given code even without an error to log", func(t *testing.T) {
		we := &logger.WideEvent{}
		we.FinishRequest(400, "VALIDATION_ERROR", nil, 0)

		fields := fixtures.Logger.WideEvent(t, we)
		assert.Equal(t, "VALIDATION_ERROR", fields["error.type"])
		assert.NotContains(t, fields, "exception.message")
	})
}

// The database cause is read from the error chain, so gateways only wrap the error with their code.
func Test_WideEvent_DatabaseCause(t *testing.T) {
	errTenantCreateFailed := errs.ErrTenantCreateFailed

	t.Run("an error the database answered logs the database and its SQLSTATE", func(t *testing.T) {
		pgErr := &pgconn.PgError{Code: "23505", Detail: "Key (slug)=(padaria) already exists."}
		we := &logger.WideEvent{}
		we.FinishRequest(500, errTenantCreateFailed.Code, fmt.Errorf("%w: %w", errTenantCreateFailed, pgErr), 0)

		fields := fixtures.Logger.WideEvent(t, we)
		assert.Equal(t, "postgresql", fields["db.system.name"])
		assert.Equal(t, "23505", fields["db.response.status_code"])
		assert.NotContains(t, fields["exception.message"], "Key (slug)", "the driver's detail may hold row values")
	})

	t.Run("a database connection failure logs the database without a SQLSTATE", func(t *testing.T) {
		we := &logger.WideEvent{}
		we.FinishCronJob(fmt.Errorf("%w: %w", errTenantCreateFailed, &pgconn.ConnectError{}), 0)

		fields := fixtures.Logger.WideEvent(t, we)
		assert.Equal(t, "postgresql", fields["db.system.name"])
		assert.NotContains(t, fields, "db.response.status_code")
	})

	t.Run("an error that did not come from the database logs no database fields", func(t *testing.T) {
		we := &logger.WideEvent{}
		we.FinishRequest(500, errs.ErrInternal.Code, errPlain, 0)

		fields := fixtures.Logger.WideEvent(t, we)
		assert.NotContains(t, fields, "db.system.name")
		assert.NotContains(t, fields, "db.response.status_code")
	})
}
