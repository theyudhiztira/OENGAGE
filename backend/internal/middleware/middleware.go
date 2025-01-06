package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"theyudhiztira/oengage-backend/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func checkCredential() {

}

type Claims struct {
	Type string `json:"type"`
	jwt.RegisteredClaims
}

var JwtSecret = []byte(config.AppEnv().JWT_SECRET)

func CheckCredential() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		username, err := ValidateToken(token)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("userId", username)
		c.Next()
	}
}
func ValidateToken(token string) (string, error) {
	const BearerSchema = "Bearer "
	if !strings.HasPrefix(token, BearerSchema) {
		return "", errors.New("Invalid token")
	}
	tokenString := token[len(BearerSchema):]
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil {
		return "", err
	}
	if !parsedToken.Valid {
		return "", errors.New("Invalid token")
	}
	if claims.Type != "access" {
		return "", errors.New("Invalid token")
	}
	return claims.Subject, nil
}
