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

// GormAdapter is the GORM implementation of the Database interface.
type GormAdapter struct {
	Conn   *gorm.DB
	config *config.DatabaseConfig
}

// New creates a new instance of GormAdapter with the given database configuration.
func New(c *config.DatabaseConfig) database.Database {
	return &GormAdapter{config: c}
}

// Migrate performs automatic migration for the database models.
func (g *GormAdapter) Migrate() error {
	return g.Conn.AutoMigrate(
		&models.User{},
	)
}

// Connect establishes a connection to the database.
func (g *GormAdapter) Connect() error {
	time.Local = time.UTC

	dsn := parseDSN(g.config)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		return err
	}

	g.Conn = db

	return nil
}

// CreateDatabase creates the database if it does not exist.
func (g *GormAdapter) CreateDatabase() error {
	var dsn = parseDSN(&config.DatabaseConfig{
		Host:         g.config.Host,
		User:         g.config.User,
		Password:     g.config.Password,
		Port:         g.config.Port,
		DatabaseName: "postgres",
	})

	conn, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}

	result := 0
	err = conn.Raw("SELECT 1 from pg_database WHERE datname=?", g.config.DatabaseName).Scan(&result).Error
	if err != nil {
		return err
	}

	hasDatabase := result > 0
	if hasDatabase {
		return nil
	}

	err = conn.Exec("CREATE DATABASE " + g.config.DatabaseName).Error
	if err != nil {
		return err
	}

	sql, err := conn.DB()
	if err != nil {
		return err
	}

	return sql.Close()
}

// Start initializes the database connection and performs migrations.
func (g *GormAdapter) Start() error {
	err := g.CreateDatabase()
	if err != nil {
		return err
	}

	err = g.Connect()
	if err != nil {
		return err
	}

	return g.Migrate()
}

// Stop closes the database connection.
func (g *GormAdapter) Stop() error {
	if g.Conn != nil {
		sqlDB, err := g.Conn.DB()
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
