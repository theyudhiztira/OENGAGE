package config

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupMongoDB() *mongo.Database {
	clientOptions := options.Client().ApplyURI(Env("MONGODB_URI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client.Database("go-api")
}
