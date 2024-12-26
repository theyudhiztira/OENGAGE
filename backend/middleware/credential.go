package middleware

import (
	"net/http"
	"strings"
	"theyudhiztira/oengage-backend/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

var JwtSecret = []byte(config.Env("JWT_SECRET"))

func CheckCredential() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		status, username := ValidateToken(token)
		if !status {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("username", username)
		c.Next()
	}
}

func ValidateToken(token string) (bool, string) {
	const BearerSchema = "Bearer "
	if !strings.HasPrefix(token, BearerSchema) {
		return false, ""
	}

	tokenString := token[len(BearerSchema):]
	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})

	if err != nil {
		return false, ""
	}

	if !parsedToken.Valid {
		return false, ""
	}

	if claims.Type != "access" {
		return false, ""
	}

	return true, claims.Username
}
