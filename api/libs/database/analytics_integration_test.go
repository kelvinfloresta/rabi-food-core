package database_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Rabi-IT/rabi-food-core/features/order/usecases"
	subscriptiongateway "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	subscriptionusecases "github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
	"github.com/Rabi-IT/rabi-food-core/fixtures"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/suite"
)

type AnalyticsSuite struct {
	suite.Suite

	app *fixtures.Api
}

func (t *AnalyticsSuite) SetupSuite() {
	t.app = fixtures.NewApi()
	t.app.Start(t.T())
}

func (t *AnalyticsSuite) SetupSubTest() {
	fixtures.CleanDatabase(t.T())
}

func (t *AnalyticsSuite) TearDownSuite() {
	t.app.Stop(t.T())
}

func TestAnalyticsSuite(t *testing.T) {
	suite.Run(t, new(AnalyticsSuite))
}

// The analytics views are what the platform owner reads in Grafana: per-tenant aggregates only.
func (t *AnalyticsSuite) Test_AnalyticsViews() {
	t.Run("should count a tenant's sales per month in Brazil's time, with the paid amount", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		ownerAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		productID := fixtures.Product.Create(t.T(), nil, ownerAuth)
		systemToken := fixtures.Auth.SystemToken(t.T())

		saoPaulo, err := time.LoadLocation("America/Sao_Paulo")
		t.Require().NoError(err)
		now := time.Now().In(saoPaulo)
		thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		lastMonth := thisMonth.AddDate(0, -1, 0)
		// 02:00 UTC on the 1st is still the previous month in Brazil (UTC-3).
		paidLastMonth := thisMonth.Add(2 * time.Hour)

		buy := func(paidAt time.Time) uint {
			customer := fixtures.Auth.UserAuth(t.T(), tenant.ID, "")
			orderID := fixtures.Order.Create(t.T(), &usecases.CreateInput{
				Items: []usecases.OrderItem{{ProductID: productID, Quantity: 1}},
			}, customer)
			order, status := fixtures.Order.GetByID(t.T(), orderID, customer)
			t.Require().Equal(http.StatusOK, status)
			t.Require().Nil(fixtures.Order.ConfirmPayment(t.T(), orderID, &usecases.ConfirmPaymentInput{
				ExternalPaymentID: orderID,
				Provider:          "stripe",
				PaidAt:            paidAt,
			}, systemToken))

			return order.TotalPrice
		}
		price := buy(paidLastMonth)
		buy(time.Now())

		// One checkout of two products is one subscription for the customer.
		fixtures.Subscription.UpsertConfig(t.T(), nil, ownerAuth)
		fixtures.Subscription.Create(t.T(), &subscriptionusecases.CreateInput{
			CheckoutKey: uuid.NewString(),
			Items: []subscriptionusecases.SubscriptionItemInput{
				{ProductID: productID, Quantity: 1},
				{ProductID: fixtures.Product.Create(t.T(), nil, ownerAuth), Quantity: 1},
			},
			DeliveryDays: []subscriptiongateway.DeliveryDay{{
				Weekday:   fixtures.Subscription.DEFAULT_DELIVERY_WEEKDAY,
				StartHour: fixtures.Subscription.DEFAULT_DELIVERY_START_HOUR,
				EndHour:   fixtures.Subscription.DEFAULT_DELIVERY_END_HOUR,
			}},
			TotalCycles: fixtures.Subscription.DEFAULT_TOTAL_CYCLES,
		}, ownerAuth)

		// Another tenant's paid sale this month never reaches this tenant's rows (new_buyers would be 2).
		other := fixtures.Tenant.Create(t.T(), nil)
		otherOrderID := fixtures.Order.Create(t.T(), nil, fixtures.Auth.TenantOwnerAuth(t.T(), other.ID, other.UserID))
		t.Require().Nil(fixtures.Order.ConfirmPayment(t.T(), otherOrderID, &usecases.ConfirmPaymentInput{
			ExternalPaymentID: otherOrderID,
			Provider:          "stripe",
			PaidAt:            time.Now(),
		}, systemToken))

		rows, err := fixtures.Analytics.Query(t.T(),
			"SELECT month, buyers, new_buyers, orders_created, orders_paid, orders_paid_cents, subscriptions_created, "+
				"subscriptions_paid, subscriptions_paid_cents FROM analytics.tenant_monthly WHERE tenant_id = $1 "+
				"ORDER BY month", tenant.ID)
		t.Require().NoError(err)
		t.Require().Len(rows, 2)

		previous, current := rows[0], rows[1]
		t.Equal(lastMonth, previous["month"], "paid at 02:00 UTC on the 1st counts in the previous month")
		t.EqualValues(0, previous["orders_created"])
		t.EqualValues(1, previous["orders_paid"])
		t.EqualValues(price, previous["orders_paid_cents"])
		t.EqualValues(1, previous["buyers"])
		t.EqualValues(1, previous["new_buyers"])

		t.Equal(thisMonth, current["month"])
		t.EqualValues(2, current["orders_created"])
		t.EqualValues(1, current["orders_paid"])
		t.EqualValues(price, current["orders_paid_cents"])
		t.EqualValues(1, current["subscriptions_created"], "two products, one checkout")
		t.EqualValues(0, current["subscriptions_paid"], "the subscription was not paid")
		t.EqualValues(0, current["subscriptions_paid_cents"])
		t.EqualValues(1, current["buyers"], "only paid sales count as buyers")
		t.EqualValues(1, current["new_buyers"], "only the second customer is new this month")
	})

	t.Run("should list the tenant with its creation date and current activity", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		ownerAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, ownerAuth)
		fixtures.Subscription.Create(t.T(), nil, ownerAuth)

		tenants, err := fixtures.Analytics.Query(t.T(),
			"SELECT name, slug, created_at, members FROM analytics.tenants WHERE tenant_id = $1", tenant.ID)
		t.Require().NoError(err)
		t.Require().Len(tenants, 1)
		t.NotNil(tenants[0]["created_at"])
		t.EqualValues(1, tenants[0]["members"])

		current, err := fixtures.Analytics.Query(t.T(),
			"SELECT active_subscriptions, buyers, last_activity_at FROM analytics.tenant_current "+
				"WHERE tenant_id = $1", tenant.ID)
		t.Require().NoError(err)
		t.Require().Len(current, 1)
		t.EqualValues(0, current[0]["active_subscriptions"], "an unpaid subscription is not active")
		t.EqualValues(0, current[0]["buyers"])
		t.NotNil(current[0]["last_activity_at"], "an unpaid subscription is still activity")
	})

	// The dashboard reads the views directly, so a view change that breaks a panel must fail here.
	t.Run("should run every platform dashboard query as the dashboards' role", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		ownerAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Order.Create(t.T(), nil, ownerAuth)

		raw, err := os.ReadFile("../../../infra/terraform/dashboards/platform.json")
		t.Require().NoError(err)

		var dashboard struct {
			Panels []struct {
				Title   string `json:"title"`
				Targets []struct {
					RawSQL string `json:"rawSql"`
				} `json:"targets"`
			} `json:"panels"`
		}
		t.Require().NoError(json.Unmarshal(raw, &dashboard))
		t.Require().NotEmpty(dashboard.Panels)

		for _, panel := range dashboard.Panels {
			for _, target := range panel.Targets {
				_, err := fixtures.Analytics.Query(t.T(), target.RawSQL)
				t.NoError(err, panel.Title)
			}
		}
	})

	t.Run("should not let the dashboards' role read tables or per-person rows", func() {
		for _, sql := range []string{
			"SELECT 1 FROM commerce.orders",
			"SELECT 1 FROM auth.users",
			"SELECT 1 FROM analytics.sales",
		} {
			_, err := fixtures.Analytics.Query(t.T(), sql)

			pgErr := new(pgconn.PgError)
			t.Require().ErrorAs(err, &pgErr, sql)
			t.Equal("42501", pgErr.Code, "%s: insufficient privilege", sql)
		}
	})
}
