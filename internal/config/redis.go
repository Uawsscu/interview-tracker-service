package config

import (
	"context"
	"interview-tracker/internal/pkg/logs"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func NewRedis() {
	addr := EnvConfig.RedisAddr
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	logs.Logger.Printf("redis| Connected successfully!")
	Rdb = rdb
}
