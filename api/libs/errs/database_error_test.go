package errs_test

import (
	"fmt"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

// The database cause must be found under the gateway's code, wherever the chain wraps it.
func Test_DatabaseErrorInfo(t *testing.T) {
	cases := map[string]struct {
		err       error
		wantState string
		wantOK    bool
	}{
		"database answered": {
			err:       fmt.Errorf("%w: insert: %w", errs.ErrTenantCreateFailed, &pgconn.PgError{Code: "23505"}),
			wantState: "23505",
			wantOK:    true,
		},
		"connection failed": {
			err:    fmt.Errorf("%w: %w", errs.ErrTenantCreateFailed, &pgconn.ConnectError{}),
			wantOK: true,
		},
		"not the database": {
			err: errs.ErrTenantNotFound,
		},
		"no error": {},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			state, ok := errs.DatabaseErrorInfo(tc.err)
			assert.Equal(t, tc.wantOK, ok)
			assert.Equal(t, tc.wantState, state)
		})
	}
}
