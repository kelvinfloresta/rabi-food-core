package fixtures

import (
	"rabi-food-core/config"
	"rabi-food-core/libs/database/gorm_adapter"
	"rabi-food-core/libs/database/gorm_adapter/models"

	"gorm.io/gorm/logger"
)

var TestDatabase = gorm_adapter.New(config.TestDatabase)

var tables = []string{
	models.User{}.TableName(),
}

func CleanDatabase() {
	gormDatabase, ok := TestDatabase.(*gorm_adapter.GormAdapter)
	if !ok {
		panic(gormDatabase)
	}

	if gormDatabase.Conn == nil {
		if err := gormDatabase.Connect(); err != nil {
			panic(err)
		}
	}

	gormDatabase.Conn.Config.Logger = logger.Default.LogMode(logger.Info)
	for _, table := range tables {
		if err := gormDatabase.Conn.Exec("TRUNCATE " + table + " CASCADE").Error; err != nil {
			panic(err)
		}
	}
}
