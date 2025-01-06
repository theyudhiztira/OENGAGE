package template

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"theyudhiztira/oengage-backend/internal/pkg"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
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

func (tr *templateRepository) GetWhatsappConfig() (pkg.SystemConfig, error) {
	cached, err := tr.Redis.Get(*tr.Ctx, pkg.SystemConfigCacheKey).Result()
	if err == redis.Nil {
		var sysConf pkg.SystemConfig
		q := tr.DB.Collection(sysConfCol).FindOne(*tr.Ctx, bson.M{})
		if q.Err() != nil {
			return pkg.SystemConfig{}, q.Err()
		}
		q.Decode(&sysConf)

		jsonData, err := json.Marshal(sysConf)
		if err != nil {
			return pkg.SystemConfig{}, err
		}

		err = tr.Redis.Set(*tr.Ctx, pkg.SystemConfigCacheKey, jsonData, 10*time.Minute).Err()
		if err != nil {
			return pkg.SystemConfig{}, err
		}
		return sysConf, nil
	} else if err != nil {
		return pkg.SystemConfig{}, err
	}

	var sysConf pkg.SystemConfig
	err = json.Unmarshal([]byte(cached), &sysConf)
	if err != nil {
		return pkg.SystemConfig{}, err
	}

	return sysConf, nil
}

func (r *templateRepository) GetWhatsappTemplate(url string, token string) (MetaTemplateResponse, error) {
	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		return MetaTemplateResponse{}, reqErr
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, respErr := client.Do(req)
	if respErr != nil {
		return MetaTemplateResponse{}, respErr
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return MetaTemplateResponse{}, readErr
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Failed to get template from %s", url)
		return MetaTemplateResponse{}, errors.New(errMsg)
	}

	var result MetaTemplateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return MetaTemplateResponse{}, err
	}

	return result, nil
}
