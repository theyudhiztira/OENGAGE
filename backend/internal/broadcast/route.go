package broadcast

import (
	"context"
	"log"
	"theyudhiztira/oengage-backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func Router(r *gin.RouterGroup, db *mongo.Database, rds *redis.Client, ctx *context.Context) *gin.RouterGroup {
	repo := NewBroadcastRepository(*db, ctx, *rds)
	service := NewBroadcastService(repo)
	handler := NewBroadcastHandler(ctx, *service)
	module := "broadcast"
	authMiddleware := middleware.NewAuthMiddleware(*db, *rds, *ctx)

	log.Println(handler)

	broadcastRouter := r.Group("/broadcast")
	{
		broadcastRouter.POST("", authMiddleware.CheckCredential(module), func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "OK",
			})
		})
	}

	return broadcastRouter
}
