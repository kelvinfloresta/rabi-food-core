package fixtures

import (
	"rabi-food-core/libs/database/gorm_adapter/models"

	"gorm.io/gorm/logger"
)

var tables = []string{
	models.User{}.TableName(),
}

func CleanDatabase() {
	if testDB.Conn == nil {
		if err := testDB.Connect(); err != nil {
			panic(err)
		}
	}

	testDB.Conn.Config.Logger = logger.Default.LogMode(logger.Info)
	for _, table := range tables {
		if err := testDB.Conn.Exec("TRUNCATE " + table + " CASCADE").Error; err != nil {
			panic(err)
		}
	}
}
