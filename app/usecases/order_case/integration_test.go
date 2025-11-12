package order_case_test

import (
	"net/http"
	"rabi-food-core/fixtures"
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/database/gateways/order_gateway"
	"rabi-food-core/libs/database/gateways/product_gateway"
	"rabi-food-core/libs/errs"
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
	t.Run("should be able to create an order successfully with valid products", func() {
		EXPECTED_PRODUCT_NAME := "Product Name"
		EXPECTED_PRODUCT_PRICE := uint(100)
		EXPECTED_TOTAL_PRICE := uint(100)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		productID := fixtures.Product.Create(t.T(), &product_gateway.CreateInput{
			Name:       "Product Name",
			CategoryID: fixtures.Category.Create(t.T(), nil, token),
			Price:      100,
			IsActive:   true,
		}, token)

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
		t.Equal(EXPECTED_PRODUCT_NAME, orderFound.Items[0].ProductName)
		t.EqualValues(EXPECTED_PRODUCT_PRICE, orderFound.Items[0].UnitPrice)
		t.EqualValues(EXPECTED_TOTAL_PRICE, orderFound.TotalPrice)
	})

	t.Run("should correctly calculate total price from multiple items", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		productID1 := fixtures.Product.Create(t.T(), nil, token)
		productID2 := fixtures.Product.Create(t.T(), nil, token)

		Body := order_case.CreateInput{
			Notes: "Notes",
			Items: []order_case.OrderItem{
				{ProductID: productID1, Quantity: 1},
				{ProductID: productID2, Quantity: 2},
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

		t.Len(orderFound.Items, 2)
		t.Equal(productID1, orderFound.Items[0].ProductID)
		t.EqualValues(300, orderFound.TotalPrice)

		t.EqualValues(1, orderFound.Items[0].Quantity)
		t.EqualValues(100, orderFound.Items[0].UnitPrice)
		t.EqualValues(100, orderFound.Items[0].Total)

		t.Equal(productID2, orderFound.Items[1].ProductID)
		t.EqualValues(2, orderFound.Items[1].Quantity)
		t.EqualValues(100, orderFound.Items[1].UnitPrice)
		t.EqualValues(200, orderFound.Items[1].Total)
	})

	t.Run("should fail when no products are found for given IDs", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := order_case.CreateInput{
			Notes: "Notes",
			Items: []order_case.OrderItem{
				{ProductID: uuid.NewString(), Quantity: 1},
			},
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusNotFound).
			JSON().IsEqual(errs.ErrProductNotFound)
	})

	t.Run("should fail when some product IDs are missing in the database", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		existingProductID := fixtures.Product.Create(t.T(), nil, token)

		Body := order_case.CreateInput{
			Notes: "Notes",
			Items: []order_case.OrderItem{
				{ProductID: existingProductID, Quantity: 1},
				{ProductID: uuid.NewString(), Quantity: 1}, // Non-existing product ID
			},
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Order.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusNotFound).
			JSON().IsEqual(errs.ErrProductNotFound)
	})

	t.Run("should ignore provided tenantID and use token tenant when user is not backoffice", func() {
		t.T().Skipf("This endpoint does not accept tenantID in the input for now")
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

	t.Run("should not be able to paginate orders from another tenant", func() {
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
	t.Run("should be able to patch the notes", func() {
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
