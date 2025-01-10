package template

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

func (r *templateRepository) GetWhatsappTemplate(url, token string) (MetaTemplateResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return MetaTemplateResponse{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return MetaTemplateResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return MetaTemplateResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return MetaTemplateResponse{}, fmt.Errorf("failed to get template from %s", url)
	}

	var result MetaTemplateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return MetaTemplateResponse{}, err
	}

	return result, nil
}

func (r *templateRepository) CreateTemplate(sysConf pkg.SystemConfig, payload WhatsappTemplate) (MetaTemplateResponse, error) {
	urlStr := fmt.Sprintf("https://graph.facebook.com/v21.0/%s/message_templates", sysConf.WhatsappWabaID)

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal request body:", err)
		return MetaTemplateResponse{}, err
	}

	req, err := http.NewRequest("POST", urlStr, io.NopCloser(strings.NewReader(string(jsonBody))))
	if err != nil {
		return MetaTemplateResponse{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sysConf.WhatsappToken))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return MetaTemplateResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return MetaTemplateResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return r.handleTemplateError(respBody)
	}

	var result MetaTemplateResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return MetaTemplateResponse{}, err
	}

	return result, nil
}

func (r *templateRepository) handleTemplateError(respBody []byte) (MetaTemplateResponse, error) {
	var errResp WhatsappTemplateErrorResp
	if err := json.Unmarshal(respBody, &errResp); err != nil {
		return MetaTemplateResponse{}, err
	}

	var metaError error
	if errResp.Error.ErrorUserTitle != "" {
		metaError = errors.New(errResp.Error.ErrorUserTitle)
	} else {
		metaError = errors.New(errResp.Error.Message)
	}

	return MetaTemplateResponse{}, metaError
}
