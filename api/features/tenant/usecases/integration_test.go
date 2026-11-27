package usecases_test

import (
	"net/http"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/fixtures"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

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

func (t *TestSuite) Test_TenantIntegration_EnrollCustomer() {
	t.Run("should enroll a customer to a tenant", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)

		customerID := fixtures.Auth.SignUp(t.T(), &fixtures.SignUpInput{
			Email:    "customer@email.com",
			Password: "password123",
			Name:     "Customer",
			Phone:    "11988888888",
			TaxID:    "12345678901",
		})
		token := fixtures.Auth.UserToken(t.T(), customerID)

		httpexpect.Default(t.T(), fixtures.ApiURL).
			Request(http.MethodPost, fixtures.Tenant.URI+tenant.ID+"/enroll-customers").
			WithHeader("Authorization", "Bearer "+token).
			Expect().Status(http.StatusNoContent)
	})

	t.Run("should be idempotent", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)

		customerID := fixtures.Auth.SignUp(t.T(), &fixtures.SignUpInput{
			Email:    "customer@email.com",
			Password: "password123",
			Name:     "Customer",
			Phone:    "11988888888",
			TaxID:    "12345678901",
		})
		token := fixtures.Auth.UserToken(t.T(), customerID)

		for range 2 {
			httpexpect.Default(t.T(), fixtures.ApiURL).
				Request(http.MethodPost, fixtures.Tenant.URI+tenant.ID+"/enroll-customers").
				WithHeader("Authorization", "Bearer "+token).
				Expect().Status(http.StatusNoContent)
		}
	})
}

func (t *TestSuite) Test_TenantIntegration_Create() {
	t.Run("should create tenant with owner as member", func() {
		tenant := fixtures.Tenant.Create(t.T(), &fixtures.TenantCreateInput{
			Name:         "My Tenant",
			UserName:     "Owner Name",
			UserPhone:    "11999999999",
			UserEmail:    "owner@email.com",
			UserPassword: "password123",
			UserTaxID:    "12345678901",
		})

		token := fixtures.Auth.TenantOwnerToken(t.T(), tenant.UserID, tenant.ID)

		httpexpect.Default(t.T(), fixtures.ApiURL).
			Request(http.MethodGet, fixtures.Tenant.URI+"me").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("X-Tenant-ID", tenant.ID).
			Expect().Status(http.StatusOK).
			JSON().Object().
			ContainsSubset(map[string]any{
				"id":   tenant.ID,
				"name": "My Tenant",
			})
	})

	t.Run("should fail if name is missing", func() {
		userID := fixtures.Auth.SignUp(t.T(), &fixtures.SignUpInput{
			Email:    "owner@email.com",
			Password: "password123",
			Name:     "Owner Name",
			Phone:    "11999999999",
			TaxID:    "12345678901",
		})
		token := fixtures.Auth.UserToken(t.T(), userID)

		response := new(middlewares.ValidationErrorResponse)
		httpexpect.Default(t.T(), fixtures.ApiURL).
			Request(http.MethodPost, "/tenant").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]string{"name": ""}).
			Expect().Status(http.StatusBadRequest).
			JSON().Decode(response)

		t.Len(response.Errors, 1)
		t.Equal("name", response.Errors[0].Field)
		t.Equal("required", response.Errors[0].Tag)
	})

	t.Run("should return 409 TENANT_SLUG_TAKEN for a name already in use", func() {
		first := fixtures.Tenant.Create(t.T(), &fixtures.TenantCreateInput{
			Name:         "Padaria Central",
			UserName:     "First Owner",
			UserPhone:    "11999999999",
			UserEmail:    "first@email.com",
			UserPassword: "password123",
			UserTaxID:    "12345678901",
		})
		userID := fixtures.Auth.SignUp(t.T(), &fixtures.SignUpInput{
			Email:    "second@email.com",
			Password: "password123",
			Name:     "Second Owner",
			Phone:    "11999999999",
			TaxID:    "12345678901",
		})
		token := fixtures.Auth.UserToken(t.T(), userID)

		response := httpexpect.Default(t.T(), fixtures.ApiURL).
			Request(http.MethodPost, "/tenant").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]string{"name": "Padaria Central"}).
			Expect().Status(http.StatusConflict)

		apiErr := new(errs.ApiError)
		response.JSON().Decode(apiErr)
		t.Equal(errs.ErrTenantSlugTaken.Code, apiErr.Code)

		// The log tells the database caused it, and why, without reading the message.
		fields := fixtures.Logger.Request(t.T(), response.Header("X-Request-ID").Raw())
		t.Equal(errs.ErrTenantSlugTaken.Code, fields["error.type"])
		t.Equal("postgresql", fields["db.system.name"])
		t.Equal("23505", fields["db.response.status_code"])

		// The existing tenant keeps its data.
		firstToken := fixtures.Auth.TenantOwnerToken(t.T(), first.UserID, first.ID)
		httpexpect.Default(t.T(), fixtures.ApiURL).
			Request(http.MethodGet, fixtures.Tenant.URI+"me").
			WithHeader("Authorization", "Bearer "+firstToken).
			Expect().Status(http.StatusOK).
			JSON().Object().
			ContainsSubset(map[string]any{"id": first.ID, "name": "Padaria Central"})
	})

	t.Run("should return 400 INVALID_BODY for a malformed body", func() {
		userID := fixtures.Auth.SignUp(t.T(), &fixtures.SignUpInput{
			Email:    "owner@email.com",
			Password: "password123",
			Name:     "Owner Name",
			Phone:    "11999999999",
			TaxID:    "12345678901",
		})

		fixtures.DefaultHTTP(t.T()).Request(http.MethodPost, "/tenant").
			WithHeader("Authorization", "Bearer "+fixtures.Auth.UserToken(t.T(), userID)).
			WithHeader("Content-Type", "application/json").
			WithText("{not json").
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object().HasValue("code", errs.ErrInvalidBody.Code)
	})
}
