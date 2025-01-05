package template

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func Router(r *gin.RouterGroup, db *mongo.Database, rds *redis.Client, ctx *context.Context) *gin.RouterGroup {
	repo := NewTemplateRepository(*db, ctx, *rds)
	service := NewTermplateService(repo)
	handler := NewTemplateHandler(*service)

	templateRouter := r.Group("/template")
	{
		templateRouter.GET("", handler.GetTemplate)
	}

	return templateRouter
}
