package fixtures

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

type analyticsFixture struct{}

var Analytics = analyticsFixture{}

// Query runs sql as analytics_reader, the role the platform's dashboards read with, so a test sees exactly
// what Grafana can see. The error is returned for tests that check what the role cannot read.
func (analyticsFixture) Query(t *testing.T, sql string, args ...any) ([]map[string]any, error) {
	t.Helper()
	ctx := context.Background()

	tx, err := testDB.Pool.Begin(ctx)
	require.NoError(t, err)
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, "SET LOCAL ROLE analytics_reader")
	require.NoError(t, err)

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToMap)
}
