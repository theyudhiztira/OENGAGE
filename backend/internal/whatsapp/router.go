package whatsapp

import (
	"context"
	"theyudhiztira/oengage-backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func Router(r *gin.RouterGroup, db *mongo.Database, rds *redis.Client, ctx *context.Context) *gin.RouterGroup {
	repo := NewWhatsappRepository(*db, ctx, *rds)
	service := NewWhatsappService(repo)
	handler := NewWhatsappHandler(ctx, *service)
	module := "whatsapp"
	authMiddleware := middleware.NewAuthMiddleware(*db, *rds, *ctx)

	whatsappRouter := r.Group("/whatsapp")
	{
		whatsappRouter.POST("/config", authMiddleware.CheckCredential(module), handler.ConfigHandler)
		// whatsappRouter.POST("", authMiddleware.CheckCredential(module), handler.CreateTemplate)
	}

	return whatsappRouter
}
