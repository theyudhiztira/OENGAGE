package route

import (
	"context"
	"theyudhiztira/oengage-backend/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Auth(r *gin.RouterGroup, db *mongo.Database, ctx *context.Context) {
	authService := service.NewAuthService(*db, ctx)

	authRoute := r.Group("/auth")
	{
		authRoute.POST("/register", authService.CreateUser)
		authRoute.POST("/login", authService.Login)
	}
}
