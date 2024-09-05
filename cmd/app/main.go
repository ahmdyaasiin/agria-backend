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
	validator := config.NewValidator()
	fiber := config.NewFiber()
	redis := config.NewRedis()

	config.App(&config.AppConfig{
		App:       fiber,
		DB:        db,
		Log:       log,
		Validator: validator,
		Redis:     redis,
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
