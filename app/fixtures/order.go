package fixtures

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/order_gateway"
	"rabi-food-core/usecases/order_case"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

type orderFixture struct {
	URI string
}

var Order = orderFixture{"/order/"}

func (orderFixture) Create(t *testing.T, input *order_case.CreateInput, token string) string {
	t.Helper()
	Body := input
	if Body == nil {
		Body = &order_case.CreateInput{
			Items: []order_case.OrderItem{
				{
					ProductID: Product.Create(t, nil, token),
					Quantity:  1,
				},
			},
			Notes: "Notes",
		}
	}

	id := ""
	httpexpect.Default(t, AppURL).
		Request(http.MethodPost, Order.URI).
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(Body).
		Expect().
		Status(http.StatusCreated).
		Body().Decode(&id)

	return id
}

func (orderFixture) GetByID(t *testing.T, id string, token string) (order_gateway.GetByIDOutput, int) {
	t.Helper()
	require.NotEmpty(t, id)

	found := order_gateway.GetByIDOutput{}

	obj := httpexpect.Default(t, AppURL).
		Request(http.MethodGet, Order.URI+id).
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK)

	response := obj.Raw()

	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found, response.StatusCode
}
