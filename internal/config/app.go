package config

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type AppConfig struct {
	App       *fiber.App
	DB        *sqlx.DB
	Log       *logrus.Logger
	Validator *validator.Validate
	Redis     *redis.Client
}

func App(c *AppConfig) {

	server := &http.Config{
		App: c.App,
	}

	server.Start()
}
