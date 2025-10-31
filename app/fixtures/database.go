package fixtures

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
	"testing"

	"github.com/stretchr/testify/require"
)

var tables = []string{
	models.User{}.TableName(),
	models.Tenant{}.TableName(),
	models.Category{}.TableName(),
	models.Product{}.TableName(),
}

func CleanDatabase(t *testing.T) {
	t.Helper()
	if testDB.Conn == nil {
		err := testDB.Connect()
		require.NoError(t, err)
	}

	for _, table := range tables {
		err := testDB.Conn.Exec("TRUNCATE " + table + " CASCADE").Error
		require.NoError(t, err)
	}
}
