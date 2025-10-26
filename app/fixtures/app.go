package fixtures

import (
	"fmt"
	"net"
	"rabi-food-core/config"
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	appHost = "localhost:" + config.TestPort
	AppURL  = "http://" + appHost
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
	t.Helper()
	wg := sync.WaitGroup{}

	dbErr := error(nil)
	httpErr := error(nil)
	wg.Go(func() {
		httpErr = a.http.Start()
		dbErr = a.database.Start()
	})

	err := waitForServer()
	require.NoError(t, err)
	wg.Wait()

	require.NoError(t, dbErr, "failed to start database")
	require.NoError(t, httpErr, "failed to start HTTP server")
}

func (a *App) Stop(t *testing.T) {
	t.Helper()
	err := a.http.Stop()
	require.NoError(t, err)

	err = a.database.Stop()
	require.NoError(t, err)
}

//nolint:mnd
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

	//nolint:err113
	return fmt.Errorf("server %s did not start within %s", appHost, timeout)
}
