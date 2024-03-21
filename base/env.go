package base

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Env struct {
	DOMAIN             string `validate:"url"`
	USER_SERVER_URL    string `validate:"url"`
	IS_PROD            bool   `validate:"boolean"`
	DATABASE_URL       string `validate:"url"`
	SECRET_API_KEY     string `validate:"required"`
	ARCGIS_TORONTO_URL string `validate:"url"`
	GLITCHTIP_DSN      string `validate:"required"`
}

func (base *Base) loadEnv() {

	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warrning .env does not exist:", err)
	}
	env := Env{
		DOMAIN:             os.Getenv("DOMAIN"),
		USER_SERVER_URL:    os.Getenv("USER_SERVER_URL"),
		IS_PROD:            strToBool(os.Getenv("IS_PROD")),
		DATABASE_URL:       os.Getenv("DATABASE_URL"),
		SECRET_API_KEY:     os.Getenv("SECRET_API_KEY"),
		ARCGIS_TORONTO_URL: os.Getenv("ARCGIS_TORONTO_URL"),
		GLITCHTIP_DSN:      os.Getenv("GLITCHTIP_DSN"),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(env)
	if err != nil {
		log.Fatal("Error .env:", err)
	}

	base.Env = &env
}

func strToBool(s string) bool {
	return s == "true"
}
