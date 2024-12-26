package route

import (
	"context"
	"theyudhiztira/oengage-backend/middleware"
	"theyudhiztira/oengage-backend/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func TemplateRoute(r *gin.RouterGroup, db *mongo.Database, ctx *context.Context) {
	templateService := service.NewTemplateService(*db, ctx)

	r.GET("/template", middleware.CheckCredential(), templateService.FetchFacebookTemplate)
}
