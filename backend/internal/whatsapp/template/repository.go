package template

import (
	"context"
	"encoding/json"
	"theyudhiztira/oengage-backend/internal/pkg"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
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

func (tr *templateRepository) GetWhatsappConfig() (pkg.SystemConfig, error) {
	cached, err := tr.Redis.Get(*tr.Ctx, pkg.SystemConfigCacheKey).Result()
	if err == redis.Nil {
		return tr.fetchAndCacheWhatsappConfig()
	} else if err != nil {
		return pkg.SystemConfig{}, err
	}

	var sysConf pkg.SystemConfig
	if err := json.Unmarshal([]byte(cached), &sysConf); err != nil {
		return pkg.SystemConfig{}, err
	}

	return sysConf, nil
}

func (tr *templateRepository) fetchAndCacheWhatsappConfig() (pkg.SystemConfig, error) {
	var sysConf pkg.SystemConfig
	q := tr.DB.Collection(sysConfCol).FindOne(*tr.Ctx, bson.M{})
	if err := q.Err(); err != nil {
		return pkg.SystemConfig{}, err
	}
	if err := q.Decode(&sysConf); err != nil {
		return pkg.SystemConfig{}, err
	}

	jsonData, err := json.Marshal(sysConf)
	if err != nil {
		return pkg.SystemConfig{}, err
	}

	if err := tr.Redis.Set(*tr.Ctx, pkg.SystemConfigCacheKey, jsonData, cacheTTL).Err(); err != nil {
		return pkg.SystemConfig{}, err
	}

	return sysConf, nil
}
