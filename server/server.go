package server

import (
	"context"
	"theyudhiztira/oengage-backend/config"

	"github.com/gin-gonic/gin"
)

func Server() *gin.Engine {
	db := config.SetupMongoDB()
	ctx := context.Background()
	server := Router(db, &ctx)

	server.Run(":8080")

	return server
}
