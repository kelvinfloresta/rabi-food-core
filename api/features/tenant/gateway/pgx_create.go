package gateway

import (
	"context"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (g *PgxTenantGatewayAdapter) Create(ctx context.Context, input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	tenantSQL, tenantArgs, err := sq.
		Insert("iam.tenants").
		Columns("id", "name", "slug", "language").
		Values(id, input.Name, input.Slug, input.Language).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("%w: build tenant insert: %w", errs.ErrTenantCreateFailed, err)
	}

	memberSQL, memberArgs, err := sq.
		Insert("iam.tenant_members").
		Columns("tenant_id", "user_id", "role").
		Values(id, input.InitialMember.UserID, input.InitialMember.Role).
		Suffix("ON CONFLICT (tenant_id, user_id) DO UPDATE SET role = EXCLUDED.role").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("%w: build member insert: %w", errs.ErrTenantMemberCreateFailed, err)
	}

	wctx := context.WithoutCancel(ctx)

	tx, err := g.DB.Pool.Begin(wctx)
	if err != nil {
		return "", fmt.Errorf("%w: begin: %w", errs.ErrTenantCreateFailed, err)
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(wctx, tenantSQL, tenantArgs...)
	if errs.IsUniqueViolation(err) {
		return "", fmt.Errorf("%w: %w", errs.ErrTenantSlugTaken, err)
	}
	if err != nil {
		return "", fmt.Errorf("%w: insert tenant: %w", errs.ErrTenantCreateFailed, err)
	}

	_, err = tx.Exec(wctx, memberSQL, memberArgs...)
	if err != nil {
		return "", fmt.Errorf("%w: insert member: %w", errs.ErrTenantMemberCreateFailed, err)
	}

	err = tx.Commit(wctx)
	if err != nil {
		return "", fmt.Errorf("%w: commit: %w", errs.ErrTenantCreateFailed, err)
	}

	return id, nil
}
