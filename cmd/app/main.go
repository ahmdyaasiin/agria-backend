package main

import (
	"fmt"
	"github.com/ahmdyaasiin/agria-backend/internal/config"
	"os"
	"strconv"
)

func init() {
	config.NewENV()
}

func main() {
	log := config.NewLogrus()
	db := config.NewSQLX()
	validator := config.NewValidator(db)
	fiber := config.NewFiber()
	redis := config.NewRedis()
	facebookOAuth := config.NewOAuthFacebook()
	googleOAuth := config.NewOAuthGoogle()

	config.App(&config.AppConfig{
		App:           fiber,
		DB:            db,
		Log:           log,
		Validator:     validator.Validate,
		Redis:         redis,
		FacebookOAuth: facebookOAuth,
		GoogleOAuth:   googleOAuth,
	})

	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		panic("failed to convert APP_PORT to int: " + err.Error())
	}

	err = fiber.Listen(fmt.Sprintf(":%d", appPort))
	if err != nil {
		panic("failed to start server: " + err.Error())
	}
}
