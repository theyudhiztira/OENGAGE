package auth

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Router(r *gin.RouterGroup, db *mongo.Database, ctx *context.Context) *gin.RouterGroup {
	repository := NewAuthRepository(*db, ctx)
	service := NewAuthService(repository)
	handler := NewAuthHandler(ctx, *service)

	authRouter := r.Group("/auth")
	{
		authRouter.POST("/setup", handler.SetupSystem)
		authRouter.POST("/login", handler.LoginHandler)
	}

	return authRouter
}
