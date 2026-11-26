package errs

import "net/http"

var (
	ErrForbidden             = newErr("FORBIDDEN", http.StatusForbidden)
	ErrInvalidCredentials    = newErr("INVALID_CREDENTIALS", http.StatusUnauthorized)
	ErrInvalidRefreshToken   = newErr("INVALID_REFRESH_TOKEN", http.StatusUnauthorized)
	ErrInvalidOTP            = newErr("INVALID_OTP", http.StatusUnauthorized)
	ErrInvalidToken          = newErr("INVALID_TOKEN", http.StatusUnauthorized)
	ErrUserAlreadyExists     = newErr("USER_ALREADY_EXISTS", http.StatusConflict)
	ErrMissingTenantID       = newErr("MISSING_TENANT_ID", http.StatusBadRequest)
	ErrScopedTokenSignFailed = newErr("SCOPED_TOKEN_SIGN_FAILED", http.StatusInternalServerError)
	// ErrSignUpSignInFailed means the account was created but the automatic sign-in after it failed. Clients must tell
	// the user the sign-up succeeded and send them to sign in, not show a sign-up failure.
	ErrSignUpSignInFailed = newErr("SIGNUP_SIGNIN_FAILED", http.StatusServiceUnavailable)

	ErrAuthUserCreateFailed   = newErr("AUTH_USER_CREATE_FAILED", http.StatusServiceUnavailable)
	ErrAuthSignInFailed       = newErr("AUTH_SIGNIN_FAILED", http.StatusServiceUnavailable)
	ErrAuthRefreshFailed      = newErr("AUTH_REFRESH_FAILED", http.StatusServiceUnavailable)
	ErrAuthSignOutFailed      = newErr("AUTH_SIGNOUT_FAILED", http.StatusServiceUnavailable)
	ErrAuthOTPSendFailed      = newErr("AUTH_OTP_SEND_FAILED", http.StatusServiceUnavailable)
	ErrAuthOTPVerifyFailed    = newErr("AUTH_OTP_VERIFY_FAILED", http.StatusServiceUnavailable)
	ErrAuthUserGetFailed      = newErr("AUTH_USER_GET_FAILED", http.StatusServiceUnavailable)
	ErrAuthUserDeleteFailed   = newErr("AUTH_USER_DELETE_FAILED", http.StatusServiceUnavailable)
	ErrAuthHealthCheckFailed  = newErr("AUTH_HEALTH_CHECK_FAILED", http.StatusServiceUnavailable)
	ErrAuthUserPatchFailed    = newErr("AUTH_USER_PATCH_FAILED", http.StatusInternalServerError)
	ErrAuthUserPaginateFailed = newErr("AUTH_USER_PAGINATE_FAILED", http.StatusInternalServerError)
	// ErrAuthTokenParseFailed means GoTrue returned an access token whose claims cannot be read.
	ErrAuthTokenParseFailed = newErr("AUTH_TOKEN_PARSE_FAILED", http.StatusInternalServerError)
)
