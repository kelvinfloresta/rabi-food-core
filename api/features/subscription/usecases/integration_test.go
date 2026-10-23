package usecases_test

import (
	"net/http"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/features/subscription"
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
	"github.com/Rabi-IT/rabi-food-core/fixtures"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite

	api *fixtures.Api
}

func (t *TestSuite) SetupSuite() {
	t.api = fixtures.NewApi()
	t.api.Start(t.T())
}

func (t *TestSuite) SetupSubTest() {
	fixtures.CleanDatabase(t.T())
}

func (t *TestSuite) TearDownSuite() {
	t.api.Stop(t.T())
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (t *TestSuite) Test_SubscriptionIntegration_Create() {
	t.Run("should be able to create", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth)

		Body := &usecases.CreateInput{
			CheckoutKey: uuid.Must(uuid.NewV7()).String(),
			Items: []usecases.SubscriptionItemInput{
				{
					ProductID: fixtures.Product.Create(t.T(), nil, tenantAuth),
					Quantity:  1,
				},
			},
			DeliveryDays: []g.DeliveryDay{
				{
					Weekday:   0,
					StartHour: 10,
					EndHour:   12,
				},
			},
			TotalCycles: 10,
			Notes:       "Notes",
		}

		var ids []string
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodPost, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			JSON().Decode(&ids)

		t.NotEmpty(ids)
	})

	t.Run("should fail if required fields are empty", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth)

		Body := map[string]any{
			"autoRenew": false,
			"notes":     "",
		}

		response := new(middlewares.ValidationErrorResponse)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodPost, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			WithJSON(Body).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Decode(&response)

		t.Len(response.Errors, 4)
		requiredFields := []string{"items", "deliveryDays", "totalCycles", "checkoutKey"}
		for i := range response.Errors {
			t.Contains(requiredFields, response.Errors[i].Field)
			t.Equal("required", response.Errors[i].Tag)
		}
	})

	t.Run("should not fail if optional fields are empty", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth)

		Body := map[string]any{
			"checkoutKey": uuid.Must(uuid.NewV7()).String(),
			"items": []usecases.SubscriptionItemInput{
				{
					ProductID: fixtures.Product.Create(t.T(), nil, tenantAuth),
					Quantity:  1,
				},
			},
			"deliveryDays": []g.DeliveryDay{
				{
					Weekday:   0,
					StartHour: 10,
					EndHour:   12,
				},
			},
			"totalCycles": 10,
		}

		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodPost, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_SubscriptionIntegration_GetByID() {
	t.Run("should be able to get by id", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth)
		subscriptionID := fixtures.Subscription.Create(t.T(), nil, tenantAuth)

		response := new(g.GetByIDOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.URI+subscriptionID).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Equal(subscriptionID, response.ID)
		t.Equal(tenant.ID, response.TenantID)
		t.Equal(tenant.UserID, response.UserID)
		t.Equal(subscription.StatusActive, response.Status)

		t.Len(response.DeliveryDays, 1)
		deliveryDay := response.DeliveryDays[0]
		t.Equal(fixtures.Subscription.DEFAULT_DELIVERY_WEEKDAY, deliveryDay.Weekday)
		t.Equal(fixtures.Subscription.DEFAULT_DELIVERY_START_HOUR, deliveryDay.StartHour)
		t.Equal(fixtures.Subscription.DEFAULT_DELIVERY_END_HOUR, deliveryDay.EndHour)

		t.NotEmpty(response.ProductID)
		t.Equal(fixtures.Subscription.DEFAULT_NOTES, response.Notes)
		t.Equal(fixtures.Subscription.DEFAULT_TOTAL_CYCLES, response.TotalCycles)
		t.Equal(fixtures.Subscription.DEFAULT_TOTAL_CYCLES, response.RemainingCycles)
		t.EqualValues(0, response.CycleDiscount)
		t.Equal(fixtures.Subscription.DEFAULT_CUTOFF_OFFSET_MINUTES, response.CutoffOffsetMinutes)
		t.True(response.AutoRenew)
	})

	t.Run("should return NotFound when get by id not found", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth)
		fixtures.Subscription.Create(t.T(), nil, tenantAuth)

		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.URI+uuid.New().String()).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_SubscriptionIntegration_Paginate() {
	t.Run("should paginate own subscriptions", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth)

		for range 3 {
			fixtures.Subscription.Create(t.T(), nil, tenantAuth)
		}

		response := new(g.PaginateOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			WithQueryObject(map[string]any{"Page": 0, "PageSize": 10}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 3)
	})

	t.Run("should not see subscriptions from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherTenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), anotherTenant.ID, anotherTenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, anotherTenantAuth)

		for range 3 {
			fixtures.Subscription.Create(t.T(), nil, anotherTenantAuth)
		}

		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)

		response := new(g.PaginateOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			WithQueryObject(map[string]any{"Page": 0, "PageSize": 10}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Empty(response.Data)
	})

	// Customers may use any tenant through X-Tenant-ID, so the header alone must never expose
	// other customers' subscriptions.
	t.Run("should not see other customers' subscriptions when using their tenant header", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth)
		fixtures.Subscription.Create(t.T(), nil, tenantAuth)

		stranger := fixtures.Auth.UserAuth(t.T(), tenant.ID, "")

		response := new(g.PaginateOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+stranger.Token).
			WithHeader("X-Tenant-ID", stranger.TenantID).
			WithQueryObject(map[string]any{"Page": 0, "PageSize": 10}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Empty(response.Data)
	})

	t.Run("should let the tenant owner see the subscriptions of every customer", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth)
		customer := fixtures.Auth.UserAuth(t.T(), tenant.ID, "")
		subscriptionID := fixtures.Subscription.Create(t.T(), nil, customer)

		response := new(g.PaginateOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+tenantAuth.Token).
			WithHeader("X-Tenant-ID", tenantAuth.TenantID).
			WithQueryObject(map[string]any{"Page": 0, "PageSize": 10}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Require().Len(response.Data, 1)
		t.Equal(subscriptionID, response.Data[0].ID)
	})
}

func (t *TestSuite) Test_SubscriptionIntegration_UpsertConfig() {
	t.Run("should let the tenant owner update the config", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)

		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth)

		config := fixtures.Subscription.GetConfig(t.T(), tenant.ID)
		t.Equal(fixtures.Subscription.DEFAULT_MAX_ATTEMPTS_PER_ORDER, config.MaxAttemptsPerOrder)
		t.Equal(fixtures.Subscription.DEFAULT_CUTOFF_OFFSET_MINUTES, config.CutoffOffsetMinutes)
		t.Len(config.DiscountRules, 1)
	})

	t.Run("should not let a customer change the config of a tenant", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth)
		before := fixtures.Subscription.GetConfig(t.T(), tenant.ID)

		stranger := fixtures.Auth.UserAuth(t.T(), tenant.ID, "")

		err := new(errs.ApiError)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodPut, fixtures.Subscription.URI+"config").
			WithHeader("Authorization", "Bearer "+stranger.Token).
			WithHeader("X-Tenant-ID", stranger.TenantID).
			WithJSON(g.UpsertConfigInput{MaxAttemptsPerOrder: 1, CutoffOffsetMinutes: 1, OrderLeadMinutes: 1}).
			Expect().
			Status(http.StatusForbidden).
			JSON().Decode(err)

		t.Equal(errs.ErrForbidden.Code, err.Code)
		t.Equal(before, fixtures.Subscription.GetConfig(t.T(), tenant.ID))
	})

	t.Run("should apply the config to the tenant of the token, ignoring the header", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)

		victim := fixtures.Tenant.Create(t.T(), nil)
		victimAuth := fixtures.Auth.TenantOwnerAuth(t.T(), victim.ID, victim.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, victimAuth)
		before := fixtures.Subscription.GetConfig(t.T(), victim.ID)

		fixtures.Subscription.UpsertConfig(t.T(), &g.UpsertConfigInput{
			MaxAttemptsPerOrder: 1,
			CutoffOffsetMinutes: 1,
			OrderLeadMinutes:    1,
		}, fixtures.RequestContext{Token: tenantAuth.Token, TenantID: victim.ID})

		t.Equal(before, fixtures.Subscription.GetConfig(t.T(), victim.ID))
		t.EqualValues(1, fixtures.Subscription.GetConfig(t.T(), tenant.ID).MaxAttemptsPerOrder)
	})
}

func (t *TestSuite) Test_SubscriptionIntegration_BackofficePaginate() {
	t.Run("should see subscriptions from all tenants", func() {
		tenant1 := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth1 := fixtures.Auth.TenantOwnerAuth(t.T(), tenant1.ID, tenant1.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth1)

		tenant2 := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth2 := fixtures.Auth.TenantOwnerAuth(t.T(), tenant2.ID, tenant2.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth2)

		for range 2 {
			fixtures.Subscription.Create(t.T(), nil, tenantAuth1)
		}
		for range 2 {
			fixtures.Subscription.Create(t.T(), nil, tenantAuth2)
		}

		backofficeToken := fixtures.Auth.BackofficeToken(t.T(), tenant1.UserID)

		response := new(g.PaginateOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.BackofficeURI).
			WithHeader("Authorization", "Bearer "+backofficeToken).
			WithQueryObject(map[string]any{"Page": 0, "PageSize": 10}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 4)
	})

	t.Run("should filter by tenantId when provided", func() {
		tenant1 := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth1 := fixtures.Auth.TenantOwnerAuth(t.T(), tenant1.ID, tenant1.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth1)

		tenant2 := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth2 := fixtures.Auth.TenantOwnerAuth(t.T(), tenant2.ID, tenant2.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth2)

		for range 2 {
			fixtures.Subscription.Create(t.T(), nil, tenantAuth1)
		}
		fixtures.Subscription.Create(t.T(), nil, tenantAuth2)

		backofficeToken := fixtures.Auth.BackofficeToken(t.T(), tenant1.UserID)

		response := new(g.PaginateOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.BackofficeURI).
			WithHeader("Authorization", "Bearer "+backofficeToken).
			WithQueryObject(map[string]any{
				"Page":     0,
				"PageSize": 10,
				"tenantId": tenant1.ID,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 2)
	})
}

func (t *TestSuite) Test_SubscriptionIntegration_Delete() {
	t.Run("should not be able to delete", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth)
		subscriptionID := fixtures.Subscription.Create(t.T(), nil, tenantAuth)

		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodDelete, fixtures.Subscription.URI+subscriptionID).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			Expect().
			Status(http.StatusMethodNotAllowed)
	})
}
