package config

import (
	"os"
)

var (
	TestPort     = testEnv("PORT", "3001")
	TestDatabase = &DatabaseConfig{
		Host:         testEnv("TEST_DATABASE_HOST", "localhost"),
		User:         testEnv("TEST_DATABASE_USER", "postgres"),
		Password:     testEnv("TEST_DATABASE_PASSWORD", "postgres"),
		DatabaseName: testEnv("TEST_DATABASE_NAME", "rabi_food_test"),
		Port:         testEnv("TEST_DATABASE_PORT", "5432"),
	}
)

func testEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
