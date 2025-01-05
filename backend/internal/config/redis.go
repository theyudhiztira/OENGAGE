package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rc := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "51I46CLtbHSkpDzL69MarN1L3n1",
	})

	ping := rc.Ping(ctx)
	if ping.Err() != nil {
		panic(ping.Err())
	}

	return rc
}
