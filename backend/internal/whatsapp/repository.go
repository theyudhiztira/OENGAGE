package whatsapp

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	whatsappConfigCollectionName = "whatsapp_configs"
	cacheTTL                     = 10 * time.Minute
)

type whatsappRepository struct {
	DB    mongo.Database
	Ctx   *context.Context
	Redis redis.Client
}

type WhatsappConfig struct {
	ID                  primitive.ObjectID    `bson:"_id,omitempty"`
	WhatsappBusinesID   string                `bson:"whatsapp_business_id"`
	WhatsappToken       string                `bson:"whatsapp_token"`
	WhatsappPhoneNumber []WhatsappPhoneNumber `bson:"whatsapp_phone_number"`
	CreatedAt           int64                 `bson:"created_at"`
	UpdatedAt           int64                 `bson:"updated_at"`
}

func NewWhatsappRepository(db mongo.Database, ctx *context.Context, redis redis.Client) *whatsappRepository {
	return &whatsappRepository{
		DB:    db,
		Ctx:   ctx,
		Redis: redis,
	}
}

func (r *whatsappRepository) FindAndUpdate(waConfig WhatsappConfig) (WhatsappConfig, error) {
	var result WhatsappConfig

	updateOptions := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{}
	update := bson.M{"$set": bson.M{
		"whatsapp_business_id":  waConfig.WhatsappBusinesID,
		"whatsapp_token":        waConfig.WhatsappToken,
		"whatsapp_phone_number": waConfig.WhatsappPhoneNumber,
		"created_at":            time.Now().Unix(),
		"updated_at":            time.Now().Unix(),
	}}

	q := r.DB.Collection(whatsappConfigCollectionName).FindOneAndUpdate(*r.Ctx, filter, update, updateOptions)
	if q.Err() == mongo.ErrNoDocuments {
		q2 := r.DB.Collection(whatsappConfigCollectionName).FindOne(*r.Ctx, filter)
		if q2.Err() != nil {
			log.Println("[WhatsappRepository.Update] Error finding document:", q2.Err())
			return result, q2.Err()
		}

		if err := q2.Decode(&result); err != nil {
			log.Println("[WhatsappRepository.Update] Error decoding document:", err)
			return result, err
		}
	} else if err := q.Decode(&result); err != nil {
		log.Println("[WhatsappRepository.Update] Error decoding updated document:", err)
		return result, err
	}

	return result, nil
}

func (r *whatsappRepository) GetWhatsappConfig() (WhatsappConfig, error) {
	var result WhatsappConfig

	c, err := r.Redis.Get(*r.Ctx, whatsappConfigCollectionName).Result()
	if err == redis.Nil {
		return r.fetchAndCacheWhatsappConfig()
	} else if err != nil {
		return result, err
	}

	if err := json.Unmarshal([]byte(c), &result); err != nil {
		return result, err
	}

	return result, nil
}

func (r *whatsappRepository) fetchAndCacheWhatsappConfig() (WhatsappConfig, error) {
	var waConf WhatsappConfig
	q := r.DB.Collection(whatsappConfigCollectionName).FindOne(*r.Ctx, bson.M{})
	if err := q.Err(); err != nil {
		return WhatsappConfig{}, err
	}
	if err := q.Decode(&waConf); err != nil {
		return WhatsappConfig{}, err
	}

	jsonData, err := json.Marshal(waConf)
	if err != nil {
		return WhatsappConfig{}, err
	}

	if err := r.Redis.Set(*r.Ctx, whatsappConfigCollectionName, jsonData, cacheTTL).Err(); err != nil {
		return WhatsappConfig{}, err
	}

	return waConf, nil
}
