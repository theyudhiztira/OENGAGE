package server

import (
	"context"
	"theyudhiztira/oengage-backend/internal/auth"
	"theyudhiztira/oengage-backend/internal/config"
	"theyudhiztira/oengage-backend/internal/pkg"
	"theyudhiztira/oengage-backend/internal/whatsapp/template"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitServer() *gin.Engine {
	env := config.AppEnv()
	srv := gin.Default()
	db := config.DBConnect().Database(env.DB_NAME)
	ctx := context.TODO()
	redis := config.InitRedis()

	_, err := pkg.Migrate(db)
	if err != nil {
		panic(err)
	}

	seedErr := pkg.SeedDB(db)
	if seedErr != nil {
		panic(seedErr)
	}

	srv.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routerV1 := srv.Group("/v1")
	{
		auth.Router(routerV1, db, &ctx)
		template.Router(routerV1, db, redis, &ctx)
	}

	srv.Run(env.OENGAGE_BACKEND_ADDRESS)

	return srv
}
