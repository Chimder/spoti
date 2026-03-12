package rediscache

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func Conn(connStr string) *redis.Client {
	opt, err := redis.ParseURL(connStr)
	if err != nil {
		panic(err)
	}

	opt.PoolSize = 100
	opt.MinIdleConns = 10
	opt.MaxIdleConns = 50
	opt.PoolTimeout = 30 * time.Second

	return redis.NewClient(opt)
}
