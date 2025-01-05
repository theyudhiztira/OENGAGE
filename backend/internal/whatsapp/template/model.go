package template

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateHandler interface {
	GetTemplate(c *gin.Context)
}

type templateHandler struct {
	Service TemplateService
}

type TemplateService interface {
}

type templateService struct {
	Repository TemplateRepository
}

type TemplateRepository interface{}

type templateRepository struct {
	DB    mongo.Database
	Ctx   *context.Context
	Redis redis.Client
}

type TemplateQueryParam struct{}

type RouterDependency struct {
	RG    *gin.RouterGroup
	DB    *mongo.Database
	Redis *redis.Client
	Ctx   context.Context
}
