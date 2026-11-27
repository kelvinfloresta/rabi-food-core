package errs

import "net/http"

var (
	ErrTenantNotFound = newErr("TENANT_NOT_FOUND", http.StatusNotFound)
	// ErrTenantSlugTaken means another tenant already has the name: the slug, which comes from it, must be unique.
	ErrTenantSlugTaken          = newErr("TENANT_SLUG_TAKEN", http.StatusConflict)
	ErrTenantCreateFailed       = newErr("TENANT_CREATE_FAILED", http.StatusInternalServerError)
	ErrTenantMemberCreateFailed = newErr("TENANT_MEMBER_CREATE_FAILED", http.StatusInternalServerError)
)
