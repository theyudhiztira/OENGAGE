package auth

import (
	"context"
	"theyudhiztira/oengage-backend/internal/pkg"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type JwtPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthHandler interface {
	LoginHandler(c *gin.Context)
}

type authHandler struct {
	Service authService
	DB      mongo.Database
	Ctx     *context.Context
}

type AuthService interface {
	SetupSystem(d SystemSetupRequest) (JwtPair, error)
	Register(d User) (JwtPair, error)
	Login(username string, password string) (JwtPair, error)
}

type authService struct {
	Repository authRepository
}

type AuthRepository interface {
	FindUserByEmail(email string) (User, error)
	CreateUser(d User) (User, error)
	GetSystemSettings() (pkg.SystemConfig, error)
	UpdateSystemSettings(d pkg.SystemConfig) error
}

type authRepository struct {
	DB  mongo.Database
	Ctx *context.Context
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	IsActive  bool               `bson:"is_active"`
	RoleID    primitive.ObjectID `bson:"role_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" default:"current_timestamp"`
}

type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Role        string             `bson:"role"`
	Permissions []Permission       `bson:"permissions"`
	CreatedBy   primitive.ObjectID `bson:"created_by"`
	CreeatedAt  time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type Permission struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty"`
	Module         string              `bson:"module"`
	PermissionRule ReadWritePermission `bson:"permission_rule"`
}

type ReadWritePermission struct {
	Read  bool `bson:"read"`
	Write bool `bson:"write"`
}

type Module struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
}

type SystemSetupRequest struct {
	AdminName     string `json:"admin_name" binding:"required"`
	AdminEmail    string `json:"admin_email" binding:"required,email"`
	AdminPassword string `json:"admin_password" binding:"required,min=8"`
}

type JwtClaims struct {
	Type string `json:"type"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
