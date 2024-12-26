package service

import (
	"context"
	"log"
	"net/http"
	"theyudhiztira/oengage-backend/config"
	"theyudhiztira/oengage-backend/dto"
	"theyudhiztira/oengage-backend/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	CreateUser(c *gin.Context)
}

type authService struct {
	DB  mongo.Database
	Ctx *context.Context
}

func NewAuthService(db mongo.Database, ctx *context.Context) *authService {
	return &authService{DB: db, Ctx: ctx}
}

var registerDTO dto.RegisterDTO
var validate = validator.New()

func (r *authService) CreateUser(c *gin.Context) {
	if err := c.ShouldBindJSON(&registerDTO); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	if err := validate.Struct(registerDTO); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.(validator.ValidationErrors)[0].Error()})
		return
	}

	checkEmail := r.DB.Collection("users").FindOne(*r.Ctx, bson.M{"email": registerDTO.Email})
	if raw, err := checkEmail.Raw(); err == nil && raw != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "Bad Request"})
		return
	}

	res, err := r.DB.Collection("users").InsertOne(*r.Ctx, bson.M{
		"name":       registerDTO.Name,
		"email":      registerDTO.Email,
		"password":   hashPassword(registerDTO.Password),
		"created_at": time.Now().Unix(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	userIdStr := res.InsertedID.(primitive.ObjectID).Hex()
	accessToken := r.GenerateAccessToken(userIdStr)
	refreshToken := r.GenerateRefreshToken(userIdStr)

	c.JSON(200, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}

func (r *authService) Login(c *gin.Context) {
	var loginDTO dto.LoginDTO

	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	if err := validate.Struct(loginDTO); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.(validator.ValidationErrors)[0].Error()})
		return
	}

	user := r.DB.Collection("users").FindOne(*r.Ctx, bson.M{"email": loginDTO.Email})
	if user.Err() != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "Incorrect email or password"})
		return
	}

	var userObj model.User
	if err := user.Decode(&userObj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if !comparePassword(userObj.Password, loginDTO.Password) {
		c.JSON(http.StatusBadGateway, gin.H{"message": "Incorrect email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  r.GenerateAccessToken(userObj.ID.Hex()),
		"refresh_token": r.GenerateRefreshToken(userObj.ID.Hex()),
	})
}

func (r *authService) GenerateAccessToken(username string) string {
	claims := jwt.MapClaims{}
	claims["sub"] = username
	claims["exp"] = time.Now().Add(time.Hour * 168).Unix()
	claims["type"] = "access"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Env("JWT_SECRET")))
	if err != nil {
		log.Fatalf("Error in generating access token: %v", err)
	}

	return tokenString
}

func (r *authService) GenerateRefreshToken(username string) string {
	claims := jwt.MapClaims{}
	claims["sub"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["type"] = "refresh"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Env("JWT_SECRET")))
	if err != nil {
		log.Fatalf("Error in generating refresh token: %v", err)
	}

	return tokenString
}

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		log.Fatalf("Error in hashing password: %v", err)
	}
	return string(hashedPassword)
}

func comparePassword(hashedPassword, password string) bool {
	log.Printf("hashedPassword: %v, password: %v", hashedPassword, password)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("Error in comparing password: ", err)
		return false
	}
	return true
}
