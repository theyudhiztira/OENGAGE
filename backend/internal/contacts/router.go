package contacts

import (
	"context"
	"theyudhiztira/oengage-backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func Router(r *gin.RouterGroup, db *mongo.Database, rds *redis.Client, ctx *context.Context) *gin.RouterGroup {
	repo := NewContactsRepository(*db, ctx, *rds)
	service := NewContactsService(repo)
	handler := NewContactsHandler(ctx, *service)
	module := "contacts"
	authMiddleware := middleware.NewAuthMiddleware(*db, *rds, *ctx)

	contactsRouter := r.Group("/contacts")
	{

		contactsRouter.POST("", authMiddleware.CheckCredential(module), handler.Create)
	}

	return contactsRouter
}
