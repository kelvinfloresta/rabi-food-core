// We want to test `classifyCode` without exporting it from the `middlewares` package.
//
//nolint:testpackage
package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errPlain = errors.New("connection reset")

type errorResponse struct {
	Code string `json:"code"`
}

// classifyCode feeds the error_code log field, so it must be the code the client receives, and it must
// never be empty for a failed request.
func Test_ClassifyCode(t *testing.T) {
	cases := map[string]struct {
		err  error
		want string
	}{
		"validation error": {validator.ValidationErrors{}, codeValidationError},
		"api error":        {errs.ErrUserAlreadyExists, errs.ErrUserAlreadyExists.Code},
		"wrapped api error": {
			fmt.Errorf("%w: %w", errs.ErrSignUpSignInFailed, errPlain),
			errs.ErrSignUpSignInFailed.Code,
		},
		"api error with a full code": {errs.ProductNotFound("p1", "p2"), errs.ErrProductNotFound.Code},
		"fiber unauthorized":         {fiber.ErrUnauthorized, errs.ErrInvalidToken.Code},
		"fiber not found":            {fiber.ErrNotFound, errs.ErrInternal.Code},
		"plain error":                {errPlain, errs.ErrInternal.Code},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := classifyCode(c.err)
			assert.Equal(t, c.want, got)

			assert.Equal(t, got, responseCode(t, c.err), "the logged code must be the code the client receives")
		})
	}
}

// responseCode returns the code ErrorHandler sends to the client for err.
func responseCode(t *testing.T, err error) string {
	t.Helper()

	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Get("/", func(*fiber.Ctx) error { return err })

	resp, testErr := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	require.NoError(t, testErr)
	defer func() { _ = resp.Body.Close() }()

	var body errorResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))

	return body.Code
}
