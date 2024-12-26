package server

import (
	"context"
	"theyudhiztira/oengage-backend/route"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Router(db *mongo.Database, ctx *context.Context) *gin.Engine {
	router := gin.Default()

	routerv1 := router.Group("/v1")
	{
		routerv1.GET("/ping", (func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		}))
		route.Auth(routerv1, db, ctx)
		route.Broadcast(routerv1, db, ctx)
		route.TemplateRoute(routerv1, db, ctx)
	}

	return router
}
