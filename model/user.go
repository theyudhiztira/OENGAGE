package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name" validate:"required"`
	Email     string             `bson:"email" validate:"required,email"`
	Password  string             `bson:"password" validate:"required"`
	CreatedAt int64              `bson:"created_at"`
	UpdatedAt int64              `bson:"updated_at" default:"time.Now().Unix()"`
}
