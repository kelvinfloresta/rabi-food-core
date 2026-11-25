package usecases_test

import (
	"context"
	"errors"
	"testing"

	g "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
	"github.com/Rabi-IT/rabi-food-core/features/auth/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errGotrueUnavailable = errors.New("gotrue unavailable")

// signUpGateway stubs only the calls SignUp makes; any other call panics on the nil embedded interface.
type signUpGateway struct {
	g.AuthGateway

	createUserErr error
	signInErr     error
}

func (s *signUpGateway) CreateUser(_ context.Context, input g.CreateUserInput) (*g.CreateUserOutput, error) {
	if s.createUserErr != nil {
		return nil, s.createUserErr
	}

	return &g.CreateUserOutput{ID: "user-id", Email: input.Email}, nil
}

func (s *signUpGateway) SignIn(context.Context, g.SignInInput) (*g.TokenOutput, error) {
	if s.signInErr != nil {
		return nil, s.signInErr
	}

	return &g.TokenOutput{AccessToken: "access-token"}, nil
}

func signUpInput() *usecases.SignUpInput {
	return &usecases.SignUpInput{
		Email:    "user@email.com",
		Password: "password123",
		Name:     "User",
		Phone:    "11988888888",
		TaxID:    "12345678901",
	}
}

func Test_SignUp(t *testing.T) {
	t.Run("should return the token when sign-in succeeds", func(t *testing.T) {
		uc := usecases.New(&signUpGateway{}, nil)

		out, err := uc.SignUp(t.Context(), signUpInput())

		require.NoError(t, err)
		assert.Equal(t, usecases.SignUpOutput{ID: "user-id", Email: "user@email.com", AccessToken: "access-token"}, out)
	})

	t.Run("should return SIGNUP_SIGNIN_FAILED when the user is created but sign-in fails", func(t *testing.T) {
		uc := usecases.New(&signUpGateway{signInErr: errGotrueUnavailable}, nil)

		_, err := uc.SignUp(t.Context(), signUpInput())

		apiErr, ok := errors.AsType[*errs.ApiError](err)
		require.True(t, ok)
		assert.Equal(t, errs.ErrSignUpSignInFailed.Code, apiErr.Code)
		assert.ErrorIs(t, err, errGotrueUnavailable, "the sign-in cause must be kept for logging")
	})

	t.Run("should take precedence over an ApiError returned by sign-in", func(t *testing.T) {
		uc := usecases.New(&signUpGateway{signInErr: errs.ErrInvalidCredentials}, nil)

		_, err := uc.SignUp(t.Context(), signUpInput())

		apiErr, ok := errors.AsType[*errs.ApiError](err)
		require.True(t, ok)
		assert.Equal(t, errs.ErrSignUpSignInFailed.Code, apiErr.Code)
	})

	t.Run("should keep the CreateUser error when the user is not created", func(t *testing.T) {
		uc := usecases.New(&signUpGateway{createUserErr: errs.ErrUserAlreadyExists}, nil)

		_, err := uc.SignUp(t.Context(), signUpInput())

		apiErr, ok := errors.AsType[*errs.ApiError](err)
		require.True(t, ok)
		assert.Equal(t, errs.ErrUserAlreadyExists.Code, apiErr.Code)
	})
}
