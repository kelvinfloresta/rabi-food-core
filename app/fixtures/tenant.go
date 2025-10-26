package fixtures

import (
	"net/http"
	"rabi-food-core/usecases/tenant_case"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

type tenantFixture struct {
	URI string
}

var Tenant = tenantFixture{"/tenant/"}

func (tenantFixture) Create(t *testing.T, input *tenant_case.CreateInput) *tenant_case.CreateOutput {
	t.Helper()
	Body := input
	if Body == nil {
		Body = &tenant_case.CreateInput{
			Name:     "Name",
			UserName: "UserName",
			Phone:    "http://example.com/photo.png",
			Email:    "email@email.com",
		}
	}

	output := &tenant_case.CreateOutput{}
	httpexpect.Default(t, AppURL).
		Request(http.MethodPost, Tenant.URI).
		WithJSON(Body).
		Expect().Status(http.StatusCreated).
		JSON().Object().
		Decode(output)

	require.NotEqual(t, &tenant_case.CreateOutput{}, output)

	return output
}
