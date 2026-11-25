package errs

import "net/http"

var (
	ErrForbidden           = newErr("FORBIDDEN", http.StatusForbidden)
	ErrInvalidCredentials  = newErr("INVALID_CREDENTIALS", http.StatusUnauthorized)
	ErrInvalidRefreshToken = newErr("INVALID_REFRESH_TOKEN", http.StatusUnauthorized)
	ErrInvalidOTP          = newErr("INVALID_OTP", http.StatusUnauthorized)
	ErrInvalidToken        = newErr("INVALID_TOKEN", http.StatusUnauthorized)
	ErrUserAlreadyExists   = newErr("USER_ALREADY_EXISTS", http.StatusConflict)
	ErrAuthServiceFailure  = newErr("AUTH_SERVICE_FAILURE", http.StatusServiceUnavailable)
	ErrMissingTenantID     = newErr("MISSING_TENANT_ID", http.StatusBadRequest)
	// ErrSignUpSignInFailed means the account was created but the automatic sign-in after it failed. Clients must tell
	// the user the sign-up succeeded and send them to sign in, not show a sign-up failure.
	ErrSignUpSignInFailed = newErr("SIGNUP_SIGNIN_FAILED", http.StatusServiceUnavailable)
)

// SignUpSignInFailed keeps the sign-in failure as the cause, so the log shows why it failed
// while the client only sees SIGNUP_SIGNIN_FAILED.
func SignUpSignInFailed(cause error) *ApiError {
	e := *ErrSignUpSignInFailed
	e.err = cause

	return &e
}
