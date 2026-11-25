package app_context

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	"github.com/Rabi-IT/rabi-food-core/domain/tenant"
)

type sessionKey string

const sessionCtxKey sessionKey = "session"

// UserSession is the authenticated user of a request. A user acts either as a customer, of any tenant,
// or as a member acting for one tenant; TenantID and TenantRole are set only in the second case
// (docs/decisions/001-shared-customer-accounts.md).
type UserSession struct {
	UserID     string
	TenantID   string
	TenantRole tenant.Role
	Name       string
	Login      string
	Role       auth.Role
}

// IsTenant reports whether the user acts as a member of a tenant. Check it rather than TenantID alone.
func (s UserSession) IsTenant() bool {
	return s.TenantRole != "" && s.TenantID != ""
}

func GetSession(ctx context.Context) UserSession {
	session, ok := ctx.Value(sessionCtxKey).(*UserSession)
	if !ok {
		return UserSession{}
	}

	return *session
}

func WithSession(ctx context.Context, s *UserSession) context.Context {
	return context.WithValue(ctx, sessionCtxKey, s)
}
