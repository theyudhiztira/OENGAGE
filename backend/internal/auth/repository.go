package auth

import (
	"context"
	"log"
	"theyudhiztira/oengage-backend/internal/pkg"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var sysConfCol = "system_configs"
var userCol = "users"

func NewAuthRepository(db mongo.Database, ctx *context.Context) *authRepository {
	return &authRepository{
		DB:  db,
		Ctx: ctx,
	}
}

func (r *authRepository) CreateUser(d User) (User, error) {
	q, err := r.DB.Collection(userCol).InsertOne(*r.Ctx, d)
	if err != nil {
		return User{}, err
	}

	insertedID := q.InsertedID.(primitive.ObjectID)
	err = r.DB.Collection(userCol).FindOne(*r.Ctx, bson.M{"_id": insertedID}).Decode(&d)
	if err != nil {
		return User{}, err
	}
	return d, nil
}

func (r *authRepository) FindUserByEmail(email string) (User, error) {
	var user User
	q := r.DB.Collection(userCol).FindOne(*r.Ctx, bson.M{
		"email": email,
	})

	if q.Err() != nil {
		return User{}, q.Err()
	}

	q.Decode(&user)

	return user, nil
}

func (r *authRepository) GetSystemSettings() (pkg.SystemConfig, error) {
	var sysConf pkg.SystemConfig
	q := r.DB.Collection(sysConfCol).FindOne(*r.Ctx, bson.M{})

	if q.Err() != nil {
		return pkg.SystemConfig{}, q.Err()
	}

	q.Decode(&sysConf)

	return sysConf, nil
}

func (r *authRepository) UpdateSystemSettings(d pkg.SystemConfig) error {
	_, err := r.DB.Collection(sysConfCol).UpdateOne(*r.Ctx, bson.D{}, bson.M{
		"$set": d,
	})

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
