package config

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler"
	"github.com/ahmdyaasiin/agria-backend/internal/repository"
	"github.com/ahmdyaasiin/agria-backend/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type AppConfig struct {
	App           *fiber.App
	DB            *sqlx.DB
	Log           *logrus.Logger
	Validator     *validator.Validate
	Redis         *redis.Client
	FacebookOAuth *oauth2.Config
	GoogleOAuth   *oauth2.Config
}

func App(c *AppConfig) {

	// repositories
	userRepository := repository.NewUserRepository(c.DB)
	addressRepository := repository.NewAddressRepository(c.DB)
	refreshRepository := repository.NewRefreshRepository(c.DB)

	// usecases
	userUseCase := usecase.NewUserUseCase(c.DB, c.Log, c.Redis, userRepository, addressRepository, refreshRepository)

	userHandler := handler.NewUserHandler(c.Log, c.Validator, c.FacebookOAuth, c.GoogleOAuth, userUseCase)

	server := &http.Config{
		App:         c.App,
		UserHandler: userHandler,
	}

	server.Start()
}
