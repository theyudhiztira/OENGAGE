package contacts

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func Router(r *gin.RouterGroup, db *mongo.Database, rds *redis.Client, ctx *context.Context) *gin.RouterGroup {
	// repo := NewContactsRepository(*db, ctx, *rds)
	// service := NewContactsService(repo)
	// handler := NewContactsHandler(ctx, *service)
	// module := "contacts"

	contactsRouter := r.Group("/contacts")
	{
		// contactsRouter.POST("/create", handler.CreateContact)
		// contactsRouter.GET("/list", handler.ListContact)
	}

	return contactsRouter
}
