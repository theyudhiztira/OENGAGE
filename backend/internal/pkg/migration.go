package pkg

import (
	"context"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Index struct {
	Key     string                 `json:"key"`
	Value   interface{}            `json:"value"`
	Options map[string]interface{} `json:"options"`
}

type CollectionIndex struct {
	CollectionName string    `json:"collectionName"`
	Indexes        [][]Index `json:"indexes"`
}

func Migrate(db *mongo.Database) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	existingCollections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		panic(err)
	}

	existingCollectionsMap := make(map[string]bool)
	for _, collection := range existingCollections {
		existingCollectionsMap[collection] = true
	}

	for _, collectionIndex := range DBCollections {
		if !existingCollectionsMap[collectionIndex.CollectionName] {
			MigrateCollection(db, collectionIndex)
			continue
		}
	}

	return true, nil
}

func MigrateCollection(db *mongo.Database, c CollectionIndex) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db.CreateCollection(ctx, c.CollectionName)

	collection := db.Collection(c.CollectionName)

	for _, index := range c.Indexes {
		keys := bson.D{}
		idOpts := options.Index()

		for _, i := range index {
			if reflect.TypeOf(i.Value).Kind() == reflect.Float64 {
				i.Value = int(i.Value.(float64))
			}

			keys = append(keys, bson.E{Key: i.Key, Value: i.Value})

			for k, v := range i.Options {
				switch k {
				case "unique":
					if unique, ok := v.(bool); ok {
						idOpts.SetUnique(unique)
					}
				}
			}
		}

		indexName, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    keys,
			Options: options.Index().SetUnique(true),
		})
		if err != nil {
			collection.Drop(ctx)
			log.Fatalf("Failed to create index (%v) on collection %s: %v", keys, c.CollectionName, err)
		}
		log.Printf("Created index: %s\n", indexName)
	}

	log.Println("Migrating collection:", c.CollectionName)

	return true, nil
}
