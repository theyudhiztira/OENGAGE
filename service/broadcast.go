package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type BroadcastService interface {
	CreateUser(c *gin.Context)
}

type broadcastService struct {
	DB  mongo.Database
	Ctx *context.Context
}

func NewBroadcastService(db mongo.Database, ctx *context.Context) *broadcastService {
	return &broadcastService{DB: db, Ctx: ctx}
}

func (r *broadcastService) CreateBroadcast(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Success"})
}
