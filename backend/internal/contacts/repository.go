package contacts

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type contactsRepository struct {
	DB    mongo.Database
	Ctx   *context.Context
	Redis redis.Client
}

type ContactVariableDB struct {
	ContactId       primitive.ObjectID `bson:"contact_id"`
	ContactVariable `bson:",inline"`
}

type getContacts struct {
	Limit int64
	Skip  int64
	Data  []ContactCard
}

func NewContactsRepository(db mongo.Database, ctx *context.Context, redis redis.Client) *contactsRepository {
	return &contactsRepository{
		DB:    db,
		Ctx:   ctx,
		Redis: redis,
	}
}

var contactsCollectionName = "contacts"
var contactVariablesCollectionName = "contact_variables"

func (r *contactsRepository) CreateContact(contacts []ContactCard) ([]ContactCard, error) {
	var result []ContactCard

	insertOptions := options.InsertMany().SetOrdered(false)
	documents := make([]interface{}, len(contacts))
	for i, contact := range contacts {
		documents[i] = contact
	}

	insert, insertError := r.DB.Collection(contactsCollectionName).InsertMany(*r.Ctx, documents, insertOptions)
	if insertError != nil {
		log.Println("[ContactsRepository.CreateContact] Error inserting contacts documents:", insertError)
		return []ContactCard{}, insertError
	}

	find, findError := r.DB.Collection(contactsCollectionName).Find(*r.Ctx, bson.M{"_id": bson.M{"$in": insert.InsertedIDs}})
	if findError != nil {
		log.Println("[ContactsRepository.CreateContact] Error finding inserted contacts documents:", findError)
		return []ContactCard{}, findError
	}

	if err := find.All(*r.Ctx, &result); err != nil {
		log.Println("[ContactsRepository.CreateContact] Error decoding inserted contacts documents:", err)
		return []ContactCard{}, err
	}

	return result, nil
}

func (r *contactsRepository) CreateVariable(variables []ContactVariableDB) ([]ContactVariableDB, error) {
	var result []ContactVariableDB

	insertOptions := options.InsertMany().SetOrdered(false)
	documents := make([]interface{}, len(variables))
	for i, variable := range variables {
		documents[i] = variable
	}

	insert, insertError := r.DB.Collection(contactVariablesCollectionName).InsertMany(*r.Ctx, documents, insertOptions)
	if insertError != nil {
		log.Println("[ContactsRepository.CreateVariable] Error inserting contact variables documents:", insertError)
		return []ContactVariableDB{}, insertError
	}

	find, findError := r.DB.Collection(contactVariablesCollectionName).Find(*r.Ctx, bson.M{"_id": bson.M{"$in": insert.InsertedIDs}})
	if findError != nil {
		log.Println("[ContactsRepository.CreateVariable] Error finding inserted contact variables documents:", findError)
		return []ContactVariableDB{}, findError
	}

	if err := find.All(*r.Ctx, &result); err != nil {
		log.Println("[ContactsRepository.CreateVariable] Error decoding inserted contact variables documents:", err)
		return []ContactVariableDB{}, err
	}

	return result, nil
}

// func (r *contactsRepository) GetContacts(limit int64, skip int64) ([]ContactCard, error) {
// 	var result []ContactCard

// 	findOptions := options.Find().SetLimit(limit).SetSkip(skip)
// 	find, findError := r.DB.Collection(contactsCollectionName).Find(*r.Ctx, bson.M{}, findOptions)

// 	return result, nil
// }
