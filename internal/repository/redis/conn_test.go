package rediscache_test

import (
	"context"
	"fmt"
	rediscache "github.com/Chimder/spoti/internal/repository/redis"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRedis(t *testing.T) {
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "redis:latest",
			ExposedPorts: []string{
				"6379/tcp",
			},
			WaitingFor: wait.ForLog("Ready to accept connections").
				WithStartupTimeout(120 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = container.Terminate(ctx)
	})

	mappedPort, err := container.MappedPort(ctx, "6379")
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	connStr := fmt.Sprintf("redis://%s:%s/0", host, mappedPort.Port())

	client := rediscache.Conn(connStr)
	require.NotNil(t, client, "redis is nil")

	status, err := client.Ping(ctx).Result()
	require.NoError(t, err)
	require.Equal(t, "PONG", status)

	key := "test-key"
	val := "hello world"

	err = client.Set(ctx, key, val, 10*time.Second).Err()
	require.NoError(t, err)

	got, err := client.Get(ctx, key).Result()
	require.NoError(t, err)
	require.Equal(t, val, got)
}
