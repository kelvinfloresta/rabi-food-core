package user_case_test

import (
	"net/http"
	"rabi-food-core/fixtures"
	"rabi-food-core/libs/database/gateways/user_gateway"
	"rabi-food-core/libs/http/fiber_adapter/middlewares"
	"rabi-food-core/usecases/tenant_case"
	"rabi-food-core/usecases/user_case"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	app *fixtures.App
}

func (t *TestSuite) SetupSuite() {
	t.app = fixtures.NewApp()
	t.app.Start(t.T())
}

func (t *TestSuite) SetupTest() {
	fixtures.CleanDatabase()
}

func (t *TestSuite) TearDownSuite() {
	t.app.Stop(t.T())
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (t *TestSuite) Test_UserIntegration_Create() {
	t.Run("should be able to create", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := user_case.CreateInput{
			Name:         "Name",
			TaxID:        "TaxID",
			City:         "City",
			State:        "State",
			Phone:        "Phone",
			ZIP:          "ZIP",
			SocialID:     "SocialID",
			Email:        "email@email.com",
			Neighborhood: "Neighborhood",
			Street:       "Street",
			Complement:   "Complement",
			Photo:        "http://example.com/photo.png",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.User.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})
	t.Run("should fail if required fields are empty", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := user_case.CreateInput{
			// Required fields left empty
			Name:  "",
			Email: "",

			// Optional fields
			TaxID:        "TaxID",
			City:         "City",
			State:        "State",
			Phone:        "Phone",
			ZIP:          "ZIP",
			SocialID:     "SocialID",
			Neighborhood: "Neighborhood",
			Street:       "Street",
			Photo:        "http://example.com/photo.png",
			Complement:   "Complement",
		}

		response := &middlewares.ValidationErrorResponse{}
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.User.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Decode(&response)

		t.Len(response.Errors, 2)
		for _, e := range response.Errors {
			switch e.Field {
			case "Name":
				t.Equal("required", e.Tag)
			case "Email":
				t.Equal("required", e.Tag)
			default:
				t.Fail("unexpected validation error field: " + e.Field)
			}
		}
	})

	t.Run("should not fail if all optional fields are not filled in", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := user_case.CreateInput{
			// Required fields
			Name:  "Name",
			Email: "email@email.com",

			// Optional fields left empty
			TenantID:     "",
			Photo:        "",
			TaxID:        "",
			City:         "",
			State:        "",
			Phone:        "",
			ZIP:          "",
			SocialID:     "",
			Neighborhood: "",
			Street:       "",
			Complement:   "",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.User.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})

	t.Run("should return an error if photo is invalid", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := user_case.CreateInput{
			Name:  "Name",
			Email: "email@email.com",
			Photo: "invalid-url",
		}

		response := &middlewares.ValidationErrorResponse{}
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.User.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Decode(response)

		t.Len(response.Errors, 1)
		t.Equal("Photo", response.Errors[0].Field)
		t.Equal("url", response.Errors[0].Tag)
	})
}

func (t *TestSuite) Test_UserIntegration_GetByID() {
	t.Run("should be able to retrieve by id", func() {
		newUser := &tenant_case.CreateInput{
			Name:     "Name",
			UserName: "UserName",
			Phone:    "Phone",
			Email:    "email@email.com",
		}

		tenant := fixtures.Tenant.Create(t.T(), newUser)

		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.User.URI+tenant.UserID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object().
			ContainsSubset(map[string]any{
				"tenantId": tenant.ID,
				"name":     newUser.UserName,
				"email":    newUser.Email,
				"phone":    newUser.Phone,
			})
	})
}

func (t *TestSuite) Test_UserIntegration_Paginate() {
	t.Run("should be able to paginate if is a backoffice user", func() {
		EXPECTED_NAME := "Name"
		tenant := fixtures.Tenant.Create(t.T(), &tenant_case.CreateInput{
			Name:     "Tenant Name",
			UserName: EXPECTED_NAME,
			Email:    "email@email.com",
			Phone:    "Phone",
		})
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		for i := 0; i < 5; i++ {
			fixtures.User.Create(t.T(), &user_case.CreateInput{
				Name:  EXPECTED_NAME,
				Email: "email@email.com",
			}, token)
		}

		backofficeToken := fixtures.Auth.BackofficeToken(t.T(), tenant.UserID)

		response := user_gateway.PaginateOutput{}
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.User.URI).
			WithHeader("Authorization", "Bearer "+backofficeToken).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 6)
		t.Equal(1, response.MaxPages)

		for _, user := range response.Data {
			t.NotEmpty(user.ID)
			t.Equal(EXPECTED_NAME, user.Name)
		}

	})

	t.Run("should not be able to paginate if is a common user", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		for i := 0; i < 5; i++ {
			fixtures.User.Create(t.T(), nil, token)
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.User.URI).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object().
			ContainsSubset(map[string]any{
				"data":     []any{},
				"maxPages": 0,
			})
	})
}
func (t *TestSuite) Test_UserIntegration_Patch() {
	t.Run("should be able to update", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := user_case.PatchValues{
			ZIP:        "NewZIP",
			Phone:      "NewPhone",
			Email:      "new-email@email.com",
			Street:     "NewStreet",
			SocialID:   "NewSocialID",
			TaxID:      "NewTaxID",
			City:       "NewCity",
			State:      "NewState",
			Complement: "NewComplement",
			Name:       "NewName",
			Photo:      "http://example.com/new-photo.png",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.User.URI+tenant.UserID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusOK).
			Body().NotEmpty()

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.User.URI+tenant.UserID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object().
			ContainsSubset(map[string]any{
				"tenantId":   tenant.ID,
				"email":      "new-email@email.com",
				"name":       "NewName",
				"zip":        "NewZIP",
				"phone":      "NewPhone",
				"street":     "NewStreet",
				"socialId":   "NewSocialID",
				"taxId":      "NewTaxID",
				"city":       "NewCity",
				"state":      "NewState",
				"complement": "NewComplement",
				"photo":      "http://example.com/new-photo.png",
			})
	})

	t.Run("should return an error if photo is invalid", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := user_case.PatchValues{
			Photo: "invalid-url",
		}

		response := &middlewares.ValidationErrorResponse{}
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.User.URI+tenant.UserID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Decode(response)

		t.Len(response.Errors, 1)
		t.Equal("Photo", response.Errors[0].Field)
		t.Equal("url", response.Errors[0].Tag)
	})
}

func (t *TestSuite) Test_UserIntegration_Delete() {
	t.Run("should be able to delete", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.User.URI+tenant.UserID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNoContent).
			Body().IsEmpty()

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.User.URI+tenant.UserID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound)
	})
}
