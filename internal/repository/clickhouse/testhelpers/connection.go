package testhelpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	migrationsDir = "../../../sql/migrations/clickhouse"
)

func SetupClickHouseContainer(t *testing.T) (driver.Conn, func()) {
	t.Helper()
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "clickhouse/clickhouse-server:latest",
			ExposedPorts: []string{
				"9000/tcp",
				"8123/tcp",
			},
			Env: map[string]string{
				"CLICKHOUSE_DB":                        "default",
				"CLICKHOUSE_USER":                      "test",
				"CLICKHOUSE_PASSWORD":                  "test",
				"CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT": "1",
			},
			WaitingFor: wait.ForHTTP("/ping").
				WithPort("8123/tcp").
				WithStatusCodeMatcher(func(status int) bool { return status == 200 }).
				WithStartupTimeout(120 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	nativePort, err := container.MappedPort(ctx, "9000")
	require.NoError(t, err)

	httpPort, err := container.MappedPort(ctx, "8123")
	require.NoError(t, err)

	conn, err := ch.Open(&ch.Options{
		Addr: []string{fmt.Sprintf("%s:%s", host, nativePort.Port())},
		Auth: ch.Auth{
			Database: "default",
			Username: "test",
			Password: "test",
		},
		DialTimeout:     10 * time.Second,
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		Debugf:          func(format string, v ...any) { t.Logf(format, v...) },
	})
	require.NoError(t, err)

	dsn := fmt.Sprintf("http://%s:%s@%s:%s?database=default", "test", "test", host, httpPort.Port())
	if err := runMigrations(dsn); err != nil {
		t.Fatalf("err run migrations: %v", err)
	}

	closeTestConn := func() {
		conn.Close()
		container.Terminate(ctx)
	}

	return conn, closeTestConn
}

func runMigrations(dsn string) error {
	db, err := goose.OpenDBWithDriver("clickhouse", dsn)
	if err != nil {
		return fmt.Errorf("goose open db: %w", err)
	}
	defer db.Close()

	goose.SetDialect("clickhouse")

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
