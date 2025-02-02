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
	CreatedBy   primitive.ObjectID `json:"created_by" bson:"created_by" binding:"omitempty"`
	CreatedAt   int64              `json:"created_at" bson:"created_at" binding:"omitempty"`
	UpdatedAt   int64              `json:"updated_at" bson:"updated_at" binding:"omitempty"`
}

type ContactVariable struct {
	Key   string `json:"key" bson:"key" binding:"required"`
	Value string `json:"value" bson:"value" binding:"required"`
}

type CreateContactRequest struct {
	ContactCard
	Variables []ContactVariable `json:"variables" binding:"omitempty"`
}

type GetContactResponse struct {
	Pagination struct {
		TotalData   int `json:"total_data"`
		TotalPage   int `json:"total_page"`
		CurrentPage int `json:"current_page"`
		Limit       int `json:"limit"`
	} `json:"pagination"`
	Data []ContactCard `json:"data"`
}
