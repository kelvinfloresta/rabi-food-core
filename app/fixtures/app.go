package fixtures

import (
	"fmt"
	"net"
	"rabi-food-core/config"
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
		http:     testHTTPServer,
		database: testDB,
	}
}

func (a *App) Start(t *testing.T) {
	go func() {
		if err := a.database.Start(); err != nil {
			panic(fmt.Sprintf("Could not start the database: %v", err))
		}

		if err := a.http.Start(); err != nil {
			panic(fmt.Sprintf("Could not start the server: %v", err))
		}
	}()

	err := waitForServer()
	require.NoError(t, err)
}

func (a *App) Stop(t *testing.T) {
	if err := a.http.Stop(); err != nil {
		require.NoError(t, err)
	}

	if err := a.database.Stop(); err != nil {
		require.NoError(t, err)
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
