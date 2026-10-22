package fixtures

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/Rabi-IT/rabi-food-core/config"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" database/sql driver used by goose
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Keep the images in sync with infra/api/docker-compose.db.yaml.
const (
	postgresImage = "postgres:18.3-bookworm"
	goTrueImage   = "supabase/gotrue:v2.189.0-rc.13"

	postgresAlias        = "postgres"
	postgresInternalPort = "5432"
	postgresUser         = "postgres"
	postgresPassword     = "postgres"
	postgresDatabase     = "rabi_food_test"

	goTruePort        = "9999/tcp"
	goTrueStartupTime = 2 * time.Minute

	// The test environment is hermetic: it never reads secrets from .env files.
	testAuthSecret = "integration-test-jwt-secret-with-at-least-32-chars"
)

// environment holds the containers backing one test package run.
type environment struct {
	network  *testcontainers.DockerNetwork
	postgres *postgres.PostgresContainer
	goTrue   testcontainers.Container
}

// startEnvironment boots Postgres and GoTrue in an isolated Docker network, applies the
// migrations and points config at them, so the DI container built afterwards uses them.
func startEnvironment(ctx context.Context) (*environment, error) {
	env := &environment{}

	nw, err := network.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("create network: %w", err)
	}
	env.network = nw

	env.postgres, err = postgres.Run(ctx, postgresImage,
		postgres.WithDatabase(postgresDatabase),
		postgres.WithUsername(postgresUser),
		postgres.WithPassword(postgresPassword),
		postgres.BasicWaitStrategies(),
		network.WithNetwork([]string{postgresAlias}, nw),
	)
	if err != nil {
		return env, fmt.Errorf("start postgres: %w", err)
	}

	dsn, err := env.postgres.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return env, fmt.Errorf("postgres connection string: %w", err)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return env, fmt.Errorf("open postgres: %w", err)
	}
	defer func() { _ = db.Close() }()

	// GoTrue only creates its tables inside an existing schema (same as infra/api/gotrue/migration.sh)
	_, err = db.ExecContext(ctx, "CREATE SCHEMA IF NOT EXISTS auth")
	if err != nil {
		return env, fmt.Errorf("create auth schema: %w", err)
	}

	env.goTrue, err = startGoTrue(ctx, nw)
	if err != nil {
		return env, err
	}

	// Application migrations reference auth.users, so they must run after GoTrue migrated its schema
	err = runMigrations(db)
	if err != nil {
		return env, err
	}

	err = configure(ctx, env)
	if err != nil {
		return env, err
	}

	return env, nil
}

func startGoTrue(ctx context.Context, nw *testcontainers.DockerNetwork) (testcontainers.Container, error) {
	databaseURL := (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(postgresUser, postgresPassword),
		Host:     net.JoinHostPort(postgresAlias, postgresInternalPort),
		Path:     postgresDatabase,
		RawQuery: "sslmode=disable&search_path=auth",
	}).String()

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        goTrueImage,
			ExposedPorts: []string{goTruePort},
			Networks:     []string{nw.Name},
			Env: map[string]string{
				"GOTRUE_API_HOST":               "0.0.0.0",
				"GOTRUE_API_PORT":               "9999",
				"API_EXTERNAL_URL":              "http://localhost:9999",
				"GOTRUE_DB_DRIVER":              "postgres",
				"GOTRUE_DB_DATABASE_URL":        databaseURL,
				"GOTRUE_DB_NAMESPACE":           "auth",
				"GOTRUE_JWT_SECRET":             testAuthSecret,
				"GOTRUE_JWT_EXP":                "3600",
				"GOTRUE_DISABLE_SIGNUP":         "true",
				"GOTRUE_SITE_URL":               "http://localhost:3000",
				"GOTRUE_MAILER_AUTOCONFIRM":     "false",
				"GOTRUE_EXTERNAL_EMAIL_ENABLED": "true",
				"GOTRUE_LOG_LEVEL":              "error",
			},
			WaitingFor: wait.ForHTTP("/health").WithPort(goTruePort).WithStartupTimeout(goTrueStartupTime),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, fmt.Errorf("start gotrue: %w", err)
	}

	return container, nil
}

func runMigrations(db *sql.DB) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("goose dialect: %w", err)
	}

	goose.SetLogger(goose.NopLogger())

	err = goose.Up(db, migrationsDir())
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

// migrationsDir resolves the migrations folder from this file, since `go test` runs each
// package with its own directory as the working directory.
func migrationsDir() string {
	_, file, _, _ := runtime.Caller(0) //nolint:dogsled // only the caller file path is needed

	return filepath.Join(filepath.Dir(file), "..", "libs", "database", "migrations")
}

// configure points the application config at the containers. It must run before the DI
// container is built, because the adapters capture config values when they are created.
func configure(ctx context.Context, env *environment) error {
	pgHost, err := env.postgres.Host(ctx)
	if err != nil {
		return fmt.Errorf("postgres host: %w", err)
	}

	pgPort, err := env.postgres.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return fmt.Errorf("postgres port: %w", err)
	}

	goTrueURL, err := env.goTrue.PortEndpoint(ctx, goTruePort, "http")
	if err != nil {
		return fmt.Errorf("gotrue endpoint: %w", err)
	}

	serviceKey, err := signServiceKey()
	if err != nil {
		return err
	}

	apiPort, err := freePort()
	if err != nil {
		return err
	}

	config.Database.Host = pgHost
	config.Database.Port = pgPort.Port()
	config.Database.User = postgresUser
	config.Database.Password = postgresPassword
	config.Database.DatabaseName = postgresDatabase
	config.GoTrueURL = goTrueURL
	config.GoTrueServiceKey = serviceKey
	config.AuthSecret = testAuthSecret
	config.ApiPort = apiPort

	apiHost = "localhost:" + apiPort
	ApiURL = "http://" + apiHost

	return nil
}

// signServiceKey issues the admin token GoTrue expects from the backend (role service_role).
func signServiceKey() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"role": "service_role"})

	signed, err := token.SignedString([]byte(testAuthSecret))
	if err != nil {
		return "", fmt.Errorf("sign service key: %w", err)
	}

	return signed, nil
}

// freePort lets each test package run its own API concurrently with the others.
func freePort() (string, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", fmt.Errorf("find free port: %w", err)
	}
	defer func() { _ = listener.Close() }()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		//nolint:err113
		return "", fmt.Errorf("unexpected listener address %v", listener.Addr())
	}

	return strconv.Itoa(addr.Port), nil
}

func (e *environment) terminate(ctx context.Context) error {
	for _, container := range []testcontainers.Container{e.goTrue, e.postgres} {
		err := testcontainers.TerminateContainer(container)
		if err != nil {
			return fmt.Errorf("terminate container: %w", err)
		}
	}

	if e.network != nil {
		err := e.network.Remove(ctx)
		if err != nil {
			return fmt.Errorf("remove network: %w", err)
		}
	}

	return nil
}
