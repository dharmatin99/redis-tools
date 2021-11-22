package lib

import (
	"github.com/go-redis/redis/v8"
)

func CreateRedisClient(addr string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})
}
