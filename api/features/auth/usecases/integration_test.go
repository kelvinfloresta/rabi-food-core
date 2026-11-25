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

func (t *TestSuite) Test_AuthIntegration_SwitchTenant() {
	t.Run("should give a member a token that acts as the tenant", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)

		token := fixtures.DefaultHTTP(t.T()).Request(http.MethodPost, "/auth/tenant/switch").
			WithHeader("Authorization", "Bearer "+fixtures.Auth.UserToken(t.T(), tenant.UserID)).
			WithJSON(map[string]string{"tenantId": tenant.ID}).
			Expect().
			Status(http.StatusOK).
			JSON().Object().Value("accessToken").String().NotEmpty().Raw()

		fixtures.DefaultHTTP(t.T()).Request(http.MethodGet, "/tenant/me").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object().HasValue("id", tenant.ID)
	})

	t.Run("should return 403 FORBIDDEN to a user who is not a member", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		outsider := fixtures.Auth.UserAuth(t.T(), tenant.ID, "")

		fixtures.DefaultHTTP(t.T()).Request(http.MethodPost, "/auth/tenant/switch").
			WithHeader("Authorization", "Bearer "+outsider.Token).
			WithJSON(map[string]string{"tenantId": tenant.ID}).
			Expect().
			Status(http.StatusForbidden).
			JSON().Object().HasValue("code", errs.ErrForbidden.Code)
	})

	t.Run("should return 400 INVALID_BODY for a malformed body", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)

		fixtures.DefaultHTTP(t.T()).Request(http.MethodPost, "/auth/tenant/switch").
			WithHeader("Authorization", "Bearer "+fixtures.Auth.UserToken(t.T(), tenant.UserID)).
			WithHeader("Content-Type", "application/json").
			WithText("{not json").
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object().HasValue("code", errs.ErrInvalidBody.Code)
	})
}

// The request log is what tells which tenant a customer was using when something failed.
func (t *TestSuite) Test_AuthIntegration_RequestLog() {
	t.Run("should log the tenant selected in the X-Tenant-ID header", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)

		requestID := fixtures.DefaultHTTP(t.T()).Request(http.MethodGet, "/category").
			WithHeader("X-Tenant-ID", tenant.ID).
			Expect().
			Status(http.StatusOK).
			Header("X-Request-ID").NotEmpty().Raw()

		fields := fixtures.Logger.Request(t.T(), requestID)
		t.Equal(tenant.ID, fields["selected_tenant_id"])
	})

	t.Run("should log the tenant a member acts for", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		owner := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)

		requestID := fixtures.DefaultHTTP(t.T()).Request(http.MethodGet, "/tenant/me").
			WithHeader("Authorization", "Bearer "+owner.Token).
			Expect().
			Status(http.StatusOK).
			Header("X-Request-ID").NotEmpty().Raw()

		fields := fixtures.Logger.Request(t.T(), requestID)
		t.Equal(tenant.ID, fields["member_tenant_id"])
	})

	t.Run("should not log a member tenant for a customer", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		customer := fixtures.Auth.UserAuth(t.T(), tenant.ID, "")

		requestID := fixtures.DefaultHTTP(t.T()).Request(http.MethodGet, "/user/me").
			WithHeader("Authorization", "Bearer "+customer.Token).
			Expect().
			Status(http.StatusOK).
			Header("X-Request-ID").NotEmpty().Raw()

		fields := fixtures.Logger.Request(t.T(), requestID)
		t.NotContains(fields, "member_tenant_id")
	})

	t.Run("should not log a selected tenant without the header", func() {
		requestID := fixtures.DefaultHTTP(t.T()).Request(http.MethodGet, "/tenant/slug/unknown").
			Expect().
			Header("X-Request-ID").NotEmpty().Raw()

		fields := fixtures.Logger.Request(t.T(), requestID)
		t.NotContains(fields, "selected_tenant_id")
	})
}

// foreignToken is a valid HS256 JWT signed with a key other than AUTH_SECRET.
const foreignToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"
