package errs_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type status string

func (s status) String() string { return string(s) }

// Every constructor must keep errors.Is working against its base error, and keep Code equal to the
// base Code: Code feeds the error_code log field, so variations (IDs, statuses) belong in FullCode.
func Test_ApiErrorConstructors(t *testing.T) {
	cases := map[string]struct {
		err  error
		base *errs.ApiError
	}{
		"ProductNotFound":          {errs.ProductNotFound("p1"), errs.ErrProductNotFound},
		"StatusNotModifiable":      {errs.StatusNotModifiable(status("DELIVERED")), errs.ErrStatusNotModifiable},
		"InvalidTranstion":         {errs.InvalidTranstion(status("A"), status("B")), errs.ErrInvalidStatusTransition},
		"DuplicateProduct":         {errs.DuplicateProduct("p1"), errs.ErrDuplicateProduct},
		"DuplicateDeliveryWeekday": {errs.DuplicateDeliveryWeekday(1), errs.ErrDuplicateDeliveryWeekday},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			require.ErrorIs(t, c.err, c.base)

			apiErr, ok := errors.AsType[*errs.ApiError](c.err)
			if assert.True(t, ok) {
				assert.Equal(t, c.base.Code, apiErr.Code)
			}
		})
	}
}

func Test_ApiError_JSON(t *testing.T) {
	t.Run("sends code and fullCode as separate fields", func(t *testing.T) {
		raw, err := json.Marshal(errs.ProductNotFound("p1", "p2"))

		require.NoError(t, err)
		assert.JSONEq(t, `{"code":"PRODUCT_NOT_FOUND","fullCode":"PRODUCT_NOT_FOUND__p1_p2","status":404}`, string(raw))
	})

	t.Run("omits fullCode when the error has no detail", func(t *testing.T) {
		raw, err := json.Marshal(errs.ErrUserAlreadyExists)

		require.NoError(t, err)
		assert.JSONEq(t, `{"code":"USER_ALREADY_EXISTS","status":409}`, string(raw))
	})
}
