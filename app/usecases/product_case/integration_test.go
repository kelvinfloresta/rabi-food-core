package product_case_test

import (
	"net/http"
	"rabi-food-core/fixtures"
	"rabi-food-core/libs/database/gateways/product_gateway"
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

func (t *TestSuite) SetupSubTest() {
	fixtures.CleanDatabase(t.T())
}

func (t *TestSuite) TearDownSuite() {
	t.app.Stop(t.T())
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (t *TestSuite) Test_ProductIntegration_Create() {

	t.Run("should be able to create", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		categoryID := fixtures.Category.Create(t.T(), nil, token)

		Body := product_gateway.CreateInput{
			Name:        "Name",
			Photo:       "http://example.com/photo.png",
			Description: "Description",
			CategoryID:  categoryID,
			Unit:        "Unit",
			Price:       100,
			IsActive:    true,
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Product.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})

	t.Run("should ignore provided tenantID and use token tenant when user is not backoffice", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		categoryID := fixtures.Category.Create(t.T(), nil, token)

		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		body := product_gateway.CreateInput{
			TenantID:    anotherTenant.ID,
			Name:        "Name",
			Description: "Description",
			CategoryID:  categoryID,
		}

		productID := fixtures.Product.Create(t.T(), &body, token)

		productFound, httpStatus := fixtures.Product.GetByID(t.T(), productID, token)
		t.Equal(http.StatusOK, httpStatus)
		t.Equal(tenant.ID, productFound.TenantID)
	})
}
