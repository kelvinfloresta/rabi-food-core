package tenant_case_test

import (
	"net/http"
	"rabi-food-core/fixtures"
	"rabi-food-core/usecases/tenant_case"
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
	t.app.Start()
}

func (t *TestSuite) SetupTest() {
	fixtures.CleanDatabase()
}

func (t *TestSuite) TearDownSuite() {
	t.app.Stop()
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (t *TestSuite) Test_Integration_Create__should_create() {
	Body := tenant_case.CreateInput{
		Name:     "Name",
		UserName: "UserName",
		Phone:    "Phone",
		Email:    "email@email.com",
	}

	var response tenant_case.CreateOutput
	httpexpect.Default(t.T(), fixtures.AppURL).
		Request(http.MethodPost, fixtures.Tenant.URI).
		WithJSON(Body).
		Expect().Status(http.StatusCreated).
		JSON().Decode(&response)

	token := fixtures.Auth.UserToken(t.T(), response.UserID)

	httpexpect.Default(t.T(), fixtures.AppURL).
		Request(http.MethodGet, fixtures.Tenant.URI+response.ID).
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK).
		JSON().Object().
		ContainsSubset(map[string]any{
			"id":   response.ID,
			"name": Body.Name,
		})
}
