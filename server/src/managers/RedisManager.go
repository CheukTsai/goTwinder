package managers

import (
	"github.com/go-redis/redis"
)

func NewRedisClientConnection(maxActive int, maxIdle int) *redis.Client{
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
        PoolSize: maxActive,
        MinIdleConns: maxIdle,
	})
	
	return client
}