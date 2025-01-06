package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBConnect() *mongo.Client {
	dbURI := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		AppEnv().DB_USER,
		AppEnv().DB_PASS,
		AppEnv().DB_HOST,
		AppEnv().DB_PORT,
		AppEnv().DB_NAME,
		AppEnv().DB_NAME,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client()
	clientOpts.SetMaxPoolSize(10)
	clientOpts.ApplyURI(dbURI)

	if gin.Mode() != gin.ReleaseMode {
		cmdMonitor := &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				log.Println("MongoDB Debug", evt.Command)
			},
		}

		clientOpts.SetMonitor(cmdMonitor)
	}

	dbClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}

	err = dbClient.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB!")

	return dbClient
}

const (
	sysConfCol   = "system_configs"
	userCol      = "users"
	roleCol      = "roles"
	moduleCol    = "modules"
	permissonCol = "permissions"
)
