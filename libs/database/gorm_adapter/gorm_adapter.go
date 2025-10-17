package gorm_adapter

import (
	"fmt"
	"rabi-food-core/config"
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/database/gorm_adapter/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormAdapter struct {
	Conn   *gorm.DB
	config *config.DatabaseConfig
}

func New(c *config.DatabaseConfig) database.Database {
	return &GormAdapter{config: c}
}

func (d *GormAdapter) Migrate() error {
	return d.Conn.AutoMigrate(
		&models.User{},
	)
}

func (d *GormAdapter) Connect() error {
	time.Local = time.UTC

	dsn := parseDSN(d.config)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		return err
	}

	d.Conn = db
	return nil
}

func (d *GormAdapter) CreateDatabase() error {
	var dsn = parseDSN(&config.DatabaseConfig{
		Host:         d.config.Host,
		User:         d.config.User,
		Password:     d.config.Password,
		Port:         d.config.Port,
		DatabaseName: "postgres",
	})

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}

	result := 0
	if err := db.Raw("SELECT 1 from pg_database WHERE datname=?", d.config.DatabaseName).Scan(&result).Error; err != nil {
		return err
	}

	hasDatabase := result > 0
	if hasDatabase {
		return nil
	}

	if err := db.Exec("CREATE DATABASE " + d.config.DatabaseName).Error; err != nil {
		return err
	}

	if db, err := db.DB(); err == nil {
		return db.Close()
	}

	return nil
}

func (g *GormAdapter) Start() error {
	if err := g.CreateDatabase(); err != nil {
		return err
	}

	if err := g.Connect(); err != nil {
		return err
	}

	return g.Migrate()
}

func (d *GormAdapter) Stop() error {
	if d.Conn != nil {
		sqlDB, err := d.Conn.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func parseDSN(d *config.DatabaseConfig) string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s",
		d.Host,
		d.User,
		d.Password,
		d.Port,
	)

	if d.DatabaseName != "" {
		return fmt.Sprintf("%s database=%s", dsn, d.DatabaseName)
	}

	return dsn
}
