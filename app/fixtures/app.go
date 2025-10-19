package fixtures

import (
	"fmt"
	"net"
	"rabi-food-core/config"
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/http"
	"rabi-food-core/libs/http/fiber_adapter"
	"rabi-food-core/libs/logger"
	"time"
)

var (
	appHost = fmt.Sprintf("localhost:%s", config.TestPort)
	AppURL  = fmt.Sprintf("http://%s", appHost)
)

type App struct {
	http     http.HTTPServer
	database database.Database
}

func NewApp() *App {
	time.Local = time.UTC

	return &App{
		http:     fiber_adapter.New(TestDatabase),
		database: TestDatabase,
	}
}

func (a *App) Start() {
	go func() {
		if err := a.database.Start(); err != nil {
			panic(err)
		}

		if err := a.http.Start(config.TestPort); err != nil {
			logger.L().Error().Msg("Could not start the server")
		}
	}()

	waitForServer()
}

func (a *App) Stop() {
	if err := a.http.Stop(); err != nil {
		logger.L().Error().Msg("Could not stop the server")
	}

	if err := a.database.Stop(); err != nil {
		logger.L().Error().Msg("Could not stop the database")
	}
}

func waitForServer() error {
	timeout := 60 * time.Second
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", appHost, 500*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("server %s did not start within %s", appHost, timeout)
}
