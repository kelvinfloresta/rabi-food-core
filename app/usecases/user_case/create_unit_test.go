package user_case_test

import (
	"rabi-food-core/fixtures"
	"rabi-food-core/fixtures/mocks"
	"rabi-food-core/libs/database/gateways/user_gateway"
	"rabi-food-core/usecases/user_case"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func makeSut(g user_gateway.UserGateway) *user_case.UserCase {
	return user_case.New(g)
}

func Test_Unit_Create__should_fail_if_required_fields_are_empty(t *testing.T) {
	DUMMY_GATEWAY := mocks.NewUserGateway(t)
	sut := makeSut(DUMMY_GATEWAY)

	_, err := sut.Create(fixtures.DUMMY_CONTEXT, &user_case.CreateInput{})
	fieldErrors := validator.ValidationErrors{}
	require.ErrorAs(t, err, &fieldErrors)
	require.Len(t, fieldErrors, 2)
	require.Equal(t, fieldErrors[0].Tag(), "required")
	require.Equal(t, fieldErrors[1].Tag(), "required")
}

func Test_Unit_Create__should_not_fail_if_all_optional_fields_are_not_filled_in(t *testing.T) {
	gateway := mocks.NewUserGateway(t)
	expectedID := "ANY_ID"
	gateway.On("Create", mock.Anything).Return(expectedID, nil)
	sut := user_case.New(gateway)

	id, err := sut.Create(fixtures.DUMMY_CONTEXT, &user_case.CreateInput{
		Name:         "Name",
		Email:        "email@email.com",
		TaxID:        "",
		City:         "",
		State:        "",
		Phone:        "",
		ZIP:          "",
		SocialID:     "",
		Neighborhood: "",
		Street:       "",
		Photo:        "",
		Complement:   "",
	})

	require.Nil(t, err)
	require.Equal(t, expectedID, id)
}
