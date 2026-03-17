package rediscache

import (
	"time"

	"github.com/redis/go-redis/extra/redisotel-native/v9"
	"github.com/redis/go-redis/v9"
)

func Conn(connStr string) *redis.Client {
	opt, err := redis.ParseURL(connStr)
	if err != nil {
		panic(err)
	}

	opt.PoolSize = 100
	opt.MinIdleConns = 15
	opt.MaxIdleConns = 50
	opt.PoolTimeout = 30 * time.Second

	obs := redisotel.GetObservabilityInstance()
	cfg := redisotel.NewConfig().
		WithEnabled(true).
		WithMetricGroups(
			redisotel.MetricGroupFlagCommand |
				redisotel.MetricGroupFlagConnectionBasic |
				redisotel.MetricGroupFlagResiliency |
				redisotel.MetricGroupFlagConnectionAdvanced,
		).
		WithHistogramBuckets([]float64{0.0001, 0.001, 0.01, 0.05, 0.1, 0.5, 1, 5, 10})

	if err := obs.Init(cfg); err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	return client
}
