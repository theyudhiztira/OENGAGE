package template

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

var sysConfCol = "system_configs"
var userCol = "users"

func NewTemplateRepository(db mongo.Database, ctx *context.Context, redis redis.Client) *templateRepository {
	return &templateRepository{
		DB:    db,
		Ctx:   ctx,
		Redis: redis,
	}
}
