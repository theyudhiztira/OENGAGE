package template

import (
	"context"
	"theyudhiztira/oengage-backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func Router(r *gin.RouterGroup, db *mongo.Database, rds *redis.Client, ctx *context.Context) *gin.RouterGroup {
	repo := NewTemplateRepository(*db, ctx, *rds)
	service := NewTemplateService(repo)
	handler := NewTemplateHandler(ctx, *service)

	templateRouter := r.Group("/template")
	{
		templateRouter.GET("", middleware.CheckCredential(), handler.GetTemplate)
	}

	return templateRouter
}
