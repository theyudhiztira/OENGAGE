package config

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Env struct {
	DB_HOST                 string
	DB_PORT                 string
	DB_USER                 string
	DB_PASS                 string
	DB_NAME                 string
	PASSWORD_SALT           string
	OENGAGE_BACKEND_ADDRESS string
	JWT_SECRET              string
	REDIS_HOST              string
	REDIS_PASS              string
}

func AppEnv() *Env {
	if gin.Mode() != gin.ReleaseMode {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	env := &Env{
		DB_HOST:                 os.Getenv("DB_HOST"),
		DB_PORT:                 os.Getenv("DB_PORT"),
		DB_USER:                 os.Getenv("DB_USER"),
		DB_PASS:                 os.Getenv("DB_PASS"),
		DB_NAME:                 os.Getenv("DB_NAME"),
		PASSWORD_SALT:           os.Getenv("PASSWORD_SALT"),
		OENGAGE_BACKEND_ADDRESS: os.Getenv("OENGAGE_BACKEND_ADDRESS"),
		JWT_SECRET:              os.Getenv("JWT_SECRET"),
		REDIS_HOST:              os.Getenv("REDIS_HOST"),
		REDIS_PASS:              os.Getenv("REDIS_PASS"),
	}

	env.SanitizedEnv()

	return env
}

func (e *Env) SanitizedEnv() {
	v := reflect.ValueOf(*e)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() == "" {
			message := fmt.Sprintf("Environment variable %s is missing", typeOfS.Field(i).Name)
			panic(message)
		}
	}
}
