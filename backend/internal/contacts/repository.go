package contacts

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type contactsRepository struct {
	DB    mongo.Database
	Ctx   *context.Context
	Redis redis.Client
}

func NewContactsRepository(db mongo.Database, ctx *context.Context, redis redis.Client) *contactsRepository {
	return &contactsRepository{
		DB:    db,
		Ctx:   ctx,
		Redis: redis,
	}
}

var contactsCollectionName = "contacts"

func (r *contactsRepository) CreateContact(contacts []ContactCard) ([]ContactCard, error) {
	var result []ContactCard

	insertOptions := options.InsertMany().SetOrdered(false)
	documents := make([]interface{}, len(contacts))
	for i, contact := range contacts {
		documents[i] = contact
	}

	insert, insertError := r.DB.Collection(contactsCollectionName).InsertMany(*r.Ctx, documents, insertOptions)
	if insertError != nil {
		log.Println("[ContactsRepository.CreateContact] Error inserting documents:", insertError)
		return []ContactCard{}, insertError
	}

	find, findError := r.DB.Collection(contactsCollectionName).Find(*r.Ctx, bson.M{"_id": bson.M{"$in": insert.InsertedIDs}})
	if findError != nil {
		log.Println("[ContactsRepository.CreateContact] Error finding inserted documents:", findError)
		return []ContactCard{}, findError
	}

	if err := find.All(*r.Ctx, &result); err != nil {
		log.Println("[ContactsRepository.CreateContact] Error decoding inserted documents:", err)
		return []ContactCard{}, err
	}

	return result, nil
}

func (r *contactsRepository) CreateVariable() error {
	return nil
}
