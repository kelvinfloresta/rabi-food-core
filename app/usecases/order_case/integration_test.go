package order_case_test

import (
	"net/http"
	"rabi-food-core/fixtures"
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/database/gateways/order_gateway"
	"rabi-food-core/libs/http/fiber_adapter/middlewares"
	"rabi-food-core/usecases/order_case"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
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

func (t *TestSuite) Test_OrderIntegration_Create() {
	t.Run("should be able to create", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		productID := fixtures.Product.Create(t.T(), nil, token)

		Body := order_case.CreateInput{
			Notes: "Notes",
			Items: []order_case.OrderItem{
				{ProductID: productID, Quantity: 1},
			},
		}

		orderID := httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty().Raw()

		orderFound, httpStatus := fixtures.Order.GetByID(t.T(), orderID, token)
		t.Equal(http.StatusOK, httpStatus)
		t.Equal("Notes", orderFound.Notes)

		t.Len(orderFound.Items, 1)
		t.Equal(productID, orderFound.Items[0].ProductID)
		t.EqualValues(1, orderFound.Items[0].Quantity)
		t.EqualValues(100, orderFound.Items[0].UnitPrice)
		t.EqualValues(100, orderFound.Items[0].Total)
	})

	t.Run("should be able to get by id", func() {
		t.T().Skipf("Skipping until OrderItem creation is implemented")
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		orderID := fixtures.Order.Create(t.T(), nil, token)

		found, status := fixtures.Order.GetByID(t.T(), orderID, token)

		t.Equal(http.StatusOK, status)
		t.Equal(orderID, found.ID)
		t.Equal("Notes", found.Notes)
	})

	t.Run("should return NotFound when get by id not found", func() {
		t.T().Skipf("Skipping until OrderItem creation is implemented")
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		_ = fixtures.Order.Create(t.T(), nil, token)

		NON_EXISTING_ID := uuid.New().String()

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Order.URI+NON_EXISTING_ID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})

	t.Run("should fail if required fields are empty", func() {
		t.T().Skipf("Skipping until required fields are defined")
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := order_gateway.CreateInput{
			// Required fields

			// Optional fields
			Notes: "",
		}

		response := new(middlewares.ValidationErrorResponse)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Decode(&response)

		t.Len(response.Errors, 1)
		t.Equal("Name", response.Errors[0].Field)
		t.Equal("required", response.Errors[0].Tag)
	})

	t.Run("should not fail if optional fields are empty", func() {
		t.T().Skipf("Skipping until optional fields are defined")
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := order_gateway.CreateInput{
			// Required fields

			// Optional fields
			Notes: "",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})

	t.Run("should ignore provided tenantID and use token tenant when user is not backoffice", func() {
		t.T().Skipf("Skipping until OrderItem creation is implemented")
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		anotherTenant := fixtures.Tenant.Create(t.T(), nil)

		Body := order_gateway.CreateInput{
			TenantID: anotherTenant.ID,
			Notes:    "Notes",
		}

		orderID := fixtures.Order.Create(t.T(), &Body, token)

		orderFound, httpStatus := fixtures.Order.GetByID(t.T(), orderID, token)
		t.Equal(http.StatusOK, httpStatus)
		t.Equal(tenant.ID, orderFound.TenantID)
	})
}

func (t *TestSuite) Test_OrderIntegration_GetByID() {
	t.Run("should be able to get by id", func() {
		t.T().Skipf("Skipping until OrderItem creation is implemented")
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		orderID := fixtures.Order.Create(t.T(), nil, token)

		found, status := fixtures.Order.GetByID(t.T(), orderID, token)

		t.Equal(http.StatusOK, status)
		t.Equal(orderID, found.ID)
		t.Equal("Notes", found.Notes)
	})

	t.Run("should return NotFound when get by id not found", func() {
		t.T().Skipf("Skipping until OrderItem creation is implemented")
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		_ = fixtures.Order.Create(t.T(), nil, token)

		NON_EXISTING_ID := uuid.New().String()

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Order.URI+NON_EXISTING_ID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})

	t.Run("should not be able to get a order from another tenant", func() {
		t.T().Skipf("Skipping until OrderItem creation is implemented")
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		anotherOrderID := fixtures.Order.Create(t.T(), nil, anotherToken)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Order.URI+anotherOrderID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_OrderIntegration_Paginate() {
	t.Run("should be able to paginate", func() {
		t.T().Skipf("Skipping until OrderItem creation is implemented")
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		for range 15 {
			fixtures.Order.Create(t.T(), nil, token)
		}

		response := new(order_gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithQueryObject(database.PaginateInput{
				Page:     0,
				PageSize: 10,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 10)
		t.Equal(2, response.MaxPages)

		for i := range 10 {
			t.NotEmpty(response.Data[i].ID)
			t.Equal("Notes", response.Data[i].Notes)
		}
	})

	t.Run("should not be able to paginate categories from another tenant", func() {
		t.T().Skipf("Skipping until OrderItem creation is implemented")
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)

		for range 5 {
			fixtures.Order.Create(t.T(), nil, anotherToken)
		}

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		response := new(order_gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithQueryObject(database.PaginateInput{
				Page:     0,
				PageSize: 10,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Empty(response.Data)
		t.Equal(0, response.MaxPages)
	})
}

func (t *TestSuite) Test_OrderIntegration_Patch() {
	t.Run("should be able to patch", func() {
		t.T().Skipf("Skipping until OrderItem creation is implemented")
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		orderID := fixtures.Order.Create(t.T(), nil, token)

		Body := order_gateway.PatchValues{
			Notes: "Updated Notes",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.Order.URI+orderID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusOK).
			Body().NotEmpty()

		found, status := fixtures.Order.GetByID(t.T(), orderID, token)

		t.Equal(http.StatusOK, status)
		t.Equal(orderID, found.ID)
		t.Equal("Updated Notes", found.Notes)
	})

	t.Run("should not be able to patch a order from another tenant", func() {
		t.T().Skipf("Skipping until OrderItem creation is implemented")
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		anotherOrderID := fixtures.Order.Create(t.T(), nil, anotherToken)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := order_gateway.PatchValues{
			Notes: "Updated Notes",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.Order.URI+anotherOrderID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_OrderIntegration_Delete() {
	t.Run("should be able to delete", func() {
		t.T().Skipf("Skipping until OrderItem creation is implemented")
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		orderID := fixtures.Order.Create(t.T(), nil, token)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.Order.URI+orderID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNoContent)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Order.URI+orderID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("should not be able to delete a order from another tenant", func() {
		t.T().Skipf("Skipping until OrderItem creation is implemented")
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		anotherOrderID := fixtures.Order.Create(t.T(), nil, anotherToken)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.Order.URI+anotherOrderID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}
