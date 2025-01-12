package auth

import (
	"errors"
	"log"
	"theyudhiztira/oengage-backend/internal/config"
	"theyudhiztira/oengage-backend/internal/pkg"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func NewAuthService(r *authRepository) *authService {
	return &authService{
		Repository: *r,
	}
}

func (s *authService) SetupSystem(d SystemSetupRequest) (JwtPair, error) {
	var payload User
	payload.Name = d.AdminEmail
	payload.Email = d.AdminEmail
	payload.Password = d.AdminPassword
	payload.IsActive = true
	payload.RoleID = pkg.GetAdminRoleObjID()
	payload.CreatedAt = time.Now()

	res, err := s.Register(payload)
	if err != nil {
		return JwtPair{}, err
	}

	return res, err
}

func (s *authService) Register(u User) (JwtPair, error) {
	u.Password = hashPassword(u.Password)

	res, err := s.Repository.CreateUser(u)
	if err != nil {
		return JwtPair{}, err
	}

	return JwtPair{
		AccessToken:  GenerateAccessToken(res),
		RefreshToken: GenerateRefreshToken(res),
	}, nil
}

func (s *authService) Login(username string, password string) (JwtPair, error) {
	res, err := s.Repository.FindUserByEmail(username)
	if err != nil {
		return JwtPair{}, err
	}

	if !comparePassword(res.Password, password) {
		return JwtPair{}, errors.New("Incorrect email or password")
	}

	return JwtPair{
		AccessToken:  GenerateAccessToken(res),
		RefreshToken: GenerateRefreshToken(res),
	}, nil
}

func GenerateAccessToken(u User) string {
	claims := jwt.MapClaims{}
	claims["sub"] = u.ID.Hex()
	claims["exp"] = time.Now().Add(time.Hour * 168).Unix()
	claims["type"] = "access"
	claims["role"] = u.RoleID.Hex()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppEnv().JWT_SECRET))
	if err != nil {
		log.Fatalf("Error in generating access token: %v", err)
	}

	return tokenString
}

func GenerateRefreshToken(u User) string {
	claims := jwt.MapClaims{}
	claims["sub"] = u.ID.Hex()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["type"] = "refresh"
	claims["role"] = u.RoleID.Hex()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppEnv().JWT_SECRET))
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
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}
