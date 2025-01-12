package broadcast

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type broadcastRepository struct {
	DB    mongo.Database
	Ctx   *context.Context
	Redis redis.Client
}

func NewBroadcastRepository(db mongo.Database, ctx *context.Context, redis redis.Client) *broadcastRepository {
	return &broadcastRepository{
		DB:    db,
		Ctx:   ctx,
		Redis: redis,
	}
}
