package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"theyudhiztira/oengage-backend/internal/auth"
	"theyudhiztira/oengage-backend/internal/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthMiddleware struct {
	DB    mongo.Database
	Redis redis.Client
	Ctx   context.Context
}

func NewAuthMiddleware(db mongo.Database, r redis.Client, ctx context.Context) *AuthMiddleware {
	return &AuthMiddleware{
		DB:    db,
		Redis: r,
		Ctx:   ctx,
	}
}

var JwtSecret = []byte(config.AppEnv().JWT_SECRET)

func (m *AuthMiddleware) CheckCredential(moduleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c.GetHeader("Authorization"))
		if token == "" {
			respondWithError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		claims, err := validateToken(token)
		if err != nil {
			respondWithError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if !m.checkPermission(claims.Role, moduleName, c.Request.Method) {
			respondWithError(c, http.StatusForbidden, "Forbidden")
			return
		}

		c.Set("userId", claims.ID)
		c.Next()
	}
}

func extractToken(header string) string {
	const BearerSchema = "Bearer "
	if !strings.HasPrefix(header, BearerSchema) {
		return ""
	}
	return header[len(BearerSchema):]
}

func validateToken(token string) (auth.JwtClaims, error) {
	claims := &auth.JwtClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil || !parsedToken.Valid || claims.Type != "access" {
		return auth.JwtClaims{}, errors.New("Invalid token")
	}
	return *claims, nil
}

func (m *AuthMiddleware) checkPermission(roleID, moduleName, method string) bool {
	role, err := m.getRoleByID(roleID)
	if err != nil {
		return false
	}

	for _, perm := range role.Permissions {
		log.Println(perm.Module, moduleName, perm.PermissionRule.Read, perm.PermissionRule.Write)
		if perm.Module == moduleName {
			if method == "GET" {
				return perm.PermissionRule.Read
			}
			return perm.PermissionRule.Write
		}
	}
	return false
}

func (m *AuthMiddleware) getRoleByID(roleID string) (auth.Role, error) {
	var role auth.Role
	rId, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return role, err
	}

	res, err := m.Redis.Get(m.Ctx, "config.role."+roleID).Result()
	if err == redis.Nil {
		err = m.DB.Collection("roles").FindOne(m.Ctx, bson.M{"_id": rId}).Decode(&role)
		if err == nil {
			roleM, err := json.Marshal(role)
			if err != nil {
				return auth.Role{}, err
			}

			_, err = m.Redis.Set(m.Ctx, "config.role."+roleID, roleM, 10*time.Minute).Result()
			if err != nil {
				return auth.Role{}, err
			}

			return role, nil
		}
	} else if err == nil {
		err = bson.UnmarshalExtJSON([]byte(res), true, &role)
	}

	uErr := json.Unmarshal([]byte(res), &role)
	if uErr != nil {
		return role, uErr
	}

	return role, err
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"message": message})
	c.Abort()
}
