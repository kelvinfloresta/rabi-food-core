package gateway_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"

	"github.com/stretchr/testify/require"
)

// Each auth service operation fails with its own code, so error.type in the log tells which one failed.
func TestGoTrueGateway_ServiceFailureCodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	gw := gateway.NewGoTrue(server.URL, "service-key", nil)
	ctx := context.Background()

	cases := map[string]struct {
		call func() error
		want *errs.ApiError
	}{
		"create user": {
			call: func() error { return errOf(gw.CreateUser(ctx, gateway.CreateUserInput{})) },
			want: errs.ErrAuthUserCreateFailed,
		},
		"sign in": {
			call: func() error { return errOf(gw.SignIn(ctx, gateway.SignInInput{})) },
			want: errs.ErrAuthSignInFailed,
		},
		"refresh": {
			call: func() error { return errOf(gw.RefreshToken(ctx, "token")) },
			want: errs.ErrAuthRefreshFailed,
		},
		"sign out": {
			call: func() error { return gw.SignOut(ctx, "token") },
			want: errs.ErrAuthSignOutFailed,
		},
		"send OTP": {
			call: func() error { return gw.SendOTP(ctx, gateway.SendOTPInput{}) },
			want: errs.ErrAuthOTPSendFailed,
		},
		"verify OTP": {
			call: func() error { return errOf(gw.VerifyOTP(ctx, gateway.VerifyOTPInput{})) },
			want: errs.ErrAuthOTPVerifyFailed,
		},
		"get user": {
			call: func() error { return errOf(gw.GetByID(ctx, "id")) },
			want: errs.ErrAuthUserGetFailed,
		},
		"delete user": {
			call: func() error { return errOf(gw.Delete(ctx, "id")) },
			want: errs.ErrAuthUserDeleteFailed,
		},
		"health": {
			call: func() error { return gw.Health(ctx) },
			want: errs.ErrAuthHealthCheckFailed,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := tc.call()

			apiErr, ok := errors.AsType[*errs.ApiError](err)
			require.True(t, ok, "expected an ApiError, got %v", err)
			require.Equal(t, tc.want.Code, apiErr.Code)
			require.Equal(t, http.StatusServiceUnavailable, apiErr.Status)
		})
	}
}

func errOf[T any](_ T, err error) error {
	return err
}
