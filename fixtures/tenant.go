package fixtures

import (
	"fmt"
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
	Body := input
	if Body == nil {
		Body = &tenant_case.CreateInput{
			Name:     "Name",
			UserName: "UserName",
			Phone:    "Phone",
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

	require.NotEqual(t, output, &tenant_case.CreateOutput{}, fmt.Sprintf("reponse: %s", output))

	return output
}
