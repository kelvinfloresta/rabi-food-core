package usecases_test

import (
	"net/http"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/fixtures"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite

	app *fixtures.Api
}

func (t *TestSuite) SetupSuite() {
	t.app = fixtures.NewApi()
	t.app.Start(t.T())
}

func (t *TestSuite) SetupSubTest() {
	fixtures.CleanDatabase(t.T())
}

func (t *TestSuite) TearDownSuite() {
	t.app.Stop(t.T())
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (t *TestSuite) Test_AuthIntegration_UserPaginate() {
	t.Run("GET /user should not be accessible by regular users", func() {
		userID := fixtures.Auth.SignUp(t.T(), &fixtures.SignUpInput{
			Email:    "user@email.com",
			Password: "password123",
			Name:     "Regular User",
			Phone:    "11988888888",
			TaxID:    "12345678901",
		})
		token := fixtures.Auth.UserToken(t.T(), userID)

		httpexpect.Default(t.T(), fixtures.ApiURL).
			Request(http.MethodGet, "/user").
			WithHeader("Authorization", "Bearer "+token).
			Expect().Status(http.StatusNotFound)
	})
}

// Every token failure on a protected route must be a 401 with the INVALID_TOKEN code. jwtware's
// default handler returns 400 for a missing token only because of a casing bug, so the API sets
// its own handler and these cases pin it down.
func (t *TestSuite) Test_AuthIntegration_ProtectedRouteToken() {
	cases := map[string]string{
		"without Authorization header":       "",
		"with an empty bearer token":         "Bearer ",
		"with another auth scheme":           "Basic dXNlcjpwYXNz",
		"with a token that is not a JWT":     "Bearer not-a-jwt",
		"with a token signed by another key": "Bearer " + foreignToken,
	}

	for name, authorization := range cases {
		t.Run("should return 401 INVALID_TOKEN "+name, func() {
			req := fixtures.DefaultHTTP(t.T()).Request(http.MethodGet, "/user/me")
			if authorization != "" {
				req = req.WithHeader("Authorization", authorization)
			}

			err := new(errs.ApiError)
			req.Expect().
				Status(http.StatusUnauthorized).
				JSON().Decode(err)

			t.Equal(errs.ErrInvalidToken.Code, err.Code)
		})
	}
}

// foreignToken is a valid HS256 JWT signed with a key other than AUTH_SECRET.
const foreignToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"
