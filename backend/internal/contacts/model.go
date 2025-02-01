package contacts

import "go.mongodb.org/mongo-driver/bson/primitive"

type ContactCard struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id" binding:"omitempty"`
	Name        string             `json:"name" bson:"name" binding:"required"`
	Phone       string             `json:"phone" bson:"phone" binding:"omitempty"`
	Email       string             `json:"email" bson:"email" binding:"omitempty"`
	Address     string             `json:"address" bson:"address" binding:"omitempty"`
	Tags        []string           `json:"tags" bson:"tags" binding:"omitempty"`
	HasWhatsapp bool               `json:"has_whatsapp" bson:"has_whatsapp" binding:"omitempty"`
	HasTelegram bool               `json:"has_telegram" bson:"has_telegram" binding:"omitempty"`
	Variables   []ContactVariable  `json:"variables" bson:"variables" binding:"omitempty"`
}

type ContactVariable struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type CreateContactRequest struct {
	ContactCard
	Variables []ContactVariable `json:"variables" binding:"omitempty"`
}
