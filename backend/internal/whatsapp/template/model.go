package template

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateHandler interface {
	GetTemplateHandler(c *gin.Context)
}

type templateHandler struct {
	Ctx     *context.Context
	Service templateService
}

type TemplateService interface {
}

type templateService struct {
	Repository templateRepository
}

type TemplateRepository interface{}

type templateRepository struct {
	DB    mongo.Database
	Ctx   *context.Context
	Redis redis.Client
}

type TemplateQueryParam struct {
	After         string `form:"after" binding:"omitempty"`
	Limit         string `form:"limit" binding:"omitempty"`
	Name          string `form:"name" binding:"omitempty"`
	NameOrContent string `form:"name_or_content" binding:"omitempty"`
	Language      string `form:"language" binding:"omitempty"`
	Category      string `form:"category" binding:"omitempty"`
	Status        string `form:"status" binding:"omitempty"`
}

type WhatsappTemplate struct {
	Name                  string      `json:"name"`
	ParameterFormat       string      `json:"parameter_format"`
	Components            []Component `json:"components"`
	Language              string      `json:"language"`
	Status                string      `json:"status"`
	Category              string      `json:"category"`
	ID                    string      `json:"id"`
	MessageSendTTLSeconds int         `json:"message_send_ttl_seconds,omitempty"`
}

type Component struct {
	Type                      string   `json:"type"`
	Format                    string   `json:"format,omitempty"`
	Text                      string   `json:"text,omitempty"`
	Example                   *Example `json:"example,omitempty"`
	Buttons                   []Button `json:"buttons,omitempty"`
	Cards                     []Card   `json:"cards,omitempty"`
	AddSecurityRecommendation bool     `json:"add_security_recommendation,omitempty"`
	CodeExpirationMinutes     int      `json:"code_expiration_minutes,omitempty"`
}

type Example struct {
	HeaderHandle []string   `json:"header_handle,omitempty"`
	BodyText     [][]string `json:"body_text,omitempty"`
}

type Button struct {
	Type        string   `json:"type"`
	Text        string   `json:"text"`
	PhoneNumber string   `json:"phone_number,omitempty"`
	URL         string   `json:"url,omitempty"`
	Example     []string `json:"example,omitempty"`
}

type Card struct {
	Components []Component `json:"components"`
}

type TemplateResponseDto struct {
	Data       interface{} `json:"data"`
	Pagination Cursor      `json:"pagination"`
}
type MetaTemplateResponse struct {
	Data   []WhatsappTemplate `json:"data"`
	Paging Paging             `json:"paging"`
}

type Paging struct {
	Cursors Cursor `json:"cursors"`
	Next    string `json:"next,omitempty"`
	Before  string `json:"before,omitempty"`
}
type Cursor struct {
	Before string `json:"before"`
	After  string `json:"after"`
}
