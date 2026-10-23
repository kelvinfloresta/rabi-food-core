package fixtures

import (
	"net/http"
	"testing"

	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	c "github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type subscriptionFixture struct {
	URI                            string
	BackofficeURI                  string
	DEFAULT_DISCOUNT_RULE          g.DiscountRule
	DEFAULT_CUTOFF_OFFSET_MINUTES  uint16
	DEFAULT_ORDER_LEAD_MINUTES     uint16
	DEFAULT_MAX_ATTEMPTS_PER_ORDER uint8

	DEFAULT_TOTAL_CYCLES        uint
	DEFAULT_NOTES               string
	DEFAULT_DELIVERY_WEEKDAY    uint8
	DEFAULT_DELIVERY_START_HOUR uint8
	DEFAULT_DELIVERY_END_HOUR   uint8
}

var (
	Subscription = subscriptionFixture{
		URI:           "/subscription/",
		BackofficeURI: "/backoffice/subscription/",
		DEFAULT_DISCOUNT_RULE: g.DiscountRule{
			CyclesThreshold: 5,  //nolint:mnd
			Discount:        10, //nolint:mnd
		},
		DEFAULT_CUTOFF_OFFSET_MINUTES:  uint16(0),
		DEFAULT_ORDER_LEAD_MINUTES:     uint16(120), //nolint:mnd
		DEFAULT_MAX_ATTEMPTS_PER_ORDER: uint8(1),

		DEFAULT_TOTAL_CYCLES:        uint(1),
		DEFAULT_NOTES:               "Notes",
		DEFAULT_DELIVERY_WEEKDAY:    uint8(0),
		DEFAULT_DELIVERY_START_HOUR: uint8(10), //nolint:mnd
		DEFAULT_DELIVERY_END_HOUR:   uint8(12), //nolint:mnd
	}
)

func (subscriptionFixture) UpsertConfig(t *testing.T, input *g.UpsertConfigInput, auth RequestContext) {
	t.Helper()
	Body := input
	if Body == nil {
		Body = &g.UpsertConfigInput{
			MaxAttemptsPerOrder: Subscription.DEFAULT_MAX_ATTEMPTS_PER_ORDER,
			DiscountRules: []g.DiscountRule{
				Subscription.DEFAULT_DISCOUNT_RULE,
			},
			CutoffOffsetMinutes: Subscription.DEFAULT_CUTOFF_OFFSET_MINUTES,
			OrderLeadMinutes:    Subscription.DEFAULT_ORDER_LEAD_MINUTES,
		}
	}

	DefaultHTTP(t).
		Request(http.MethodPut, Subscription.URI+"config").
		WithHeader("Authorization", "Bearer "+auth.Token).
		WithHeader("X-Tenant-ID", auth.TenantID).
		WithJSON(Body).
		Expect().
		StatusList(http.StatusOK, http.StatusCreated)
}

// GetConfig reads the public subscription config of a tenant.
func (subscriptionFixture) GetConfig(t *testing.T, tenantID string) g.GetConfigOutput {
	t.Helper()

	config := g.GetConfigOutput{}
	DefaultHTTP(t).
		Request(http.MethodGet, Subscription.URI+"config").
		WithHeader("X-Tenant-ID", tenantID).
		Expect().
		Status(http.StatusOK).
		JSON().Decode(&config)

	return config
}

// Create subscribes as auth (owner or customer). A nil input subscribes to a new product of
// auth's tenant with the default delivery and cycles.
func (subscriptionFixture) Create(t *testing.T, input *c.CreateInput, auth RequestContext) string {
	t.Helper()

	if input == nil {
		// Only tenant owners create products, so the product is created by an owner of auth's tenant.
		owner := Auth.TenantOwnerAuth(t, auth.TenantID, "")

		input = &c.CreateInput{
			CheckoutKey: uuid.Must(uuid.NewV7()).String(),
			Items: []c.SubscriptionItemInput{
				{
					ProductID: Product.Create(t, nil, owner),
					Quantity:  1,
				},
			},
			DeliveryDays: []g.DeliveryDay{
				{
					Weekday:   Subscription.DEFAULT_DELIVERY_WEEKDAY,
					StartHour: Subscription.DEFAULT_DELIVERY_START_HOUR,
					EndHour:   Subscription.DEFAULT_DELIVERY_END_HOUR,
				},
			},
			TotalCycles: Subscription.DEFAULT_TOTAL_CYCLES,
			Notes:       Subscription.DEFAULT_NOTES,
		}
	}

	var ids []string
	DefaultHTTP(t).
		Request(http.MethodPost, Subscription.URI).
		WithHeader("Authorization", "Bearer "+auth.Token).
		WithHeader("X-Tenant-ID", auth.TenantID).
		WithJSON(input).
		Expect().
		Status(http.StatusCreated).
		JSON().Decode(&ids)

	require.NotEmpty(t, ids)

	return ids[0]
}

func (subscriptionFixture) GetByID(t *testing.T, id string, auth RequestContext) (g.GetByIDOutput, int) {
	t.Helper()
	found := g.GetByIDOutput{}

	obj := DefaultHTTP(t).
		Request(http.MethodGet, Subscription.URI+id).
		WithHeader("Authorization", "Bearer "+auth.Token).
		WithHeader("X-Tenant-ID", auth.TenantID).
		Expect().Status(http.StatusOK)

	response := obj.Raw()

	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found, response.StatusCode
}
