package template

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	sysConfCol = "system_configs"
	userCol    = "users"
	cacheTTL   = 10 * time.Minute
)

func NewTemplateRepository(db mongo.Database, ctx *context.Context, redis redis.Client) *templateRepository {
	return &templateRepository{
		DB:    db,
		Ctx:   ctx,
		Redis: redis,
	}
}
