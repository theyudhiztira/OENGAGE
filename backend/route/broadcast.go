package route

import (
	"context"
	"theyudhiztira/oengage-backend/middleware"
	"theyudhiztira/oengage-backend/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Broadcast(r *gin.RouterGroup, db *mongo.Database, ctx *context.Context) {
	broadcastService := service.NewBroadcastService(*db, ctx)

	r.POST("/broadcast", middleware.CheckCredential(), broadcastService.CreateBroadcast)
}
