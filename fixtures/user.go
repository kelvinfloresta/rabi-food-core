package fixtures

import (
	"fmt"
	"net/http"
	"rabi-food-core/libs/database/gateways/user_gateway"
	"rabi-food-core/usecases/user_case"
	"testing"

	"github.com/stretchr/testify/require"
)

type userFixture struct {
	URI string
}

var User = userFixture{"/user/"}

func (userFixture) Create(t *testing.T, input *user_case.CreateInput, token string) string {
	Body := input
	if Body == nil {
		Body = &user_case.CreateInput{
			Name:         "Name",
			Photo:        "Photo",
			TaxID:        "TaxID",
			City:         "City",
			State:        "State",
			Phone:        "Phone",
			ZIP:          "ZIP",
			SocialID:     "SocialID",
			Email:        "email@email.com",
			Neighborhood: "Neighborhood",
			Street:       "Street",
			Complement:   "Complement",
		}
	}

	id := ""
	statusCode := Post(t, PostInput{
		Body:     Body,
		URI:      User.URI,
		Response: &id,
		Token:    token,
	})

	require.Equal(t, http.StatusCreated, statusCode, fmt.Sprintf("reponse: %s", id))
	require.NotEmpty(t, id)

	return id
}

func (userFixture) GetByID(t *testing.T, id string, token string) (user_gateway.GetByIDOutput, int) {
	found := user_gateway.GetByIDOutput{}

	input := GetInput{
		URI:      User.URI + id,
		Response: &found,
		Token:    token,
	}

	statusCode := Get(t, input)

	return found, statusCode
}
