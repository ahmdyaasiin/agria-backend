package config

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/middleware"
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
	productRepository := repository.NewProductRepository(c.DB)
	productMediaRepository := repository.NewProductMediaRepository(c.DB)
	ratingRepository := repository.NewRatingRepository(c.DB)
	cartRepository := repository.NewCartRepository(c.DB)
	wishlistRepository := repository.NewWishlistRepository(c.DB)

	// usecases
	userUseCase := usecase.NewUserUseCase(c.DB, c.Log, c.Redis, userRepository, addressRepository, refreshRepository)
	productUseCase := usecase.NewProductUseCase(c.DB, c.Log, c.Redis, addressRepository, productRepository, productMediaRepository, ratingRepository)
	cartUseCase := usecase.NewCartUseCase(c.DB, c.Log, c.Redis, cartRepository, productRepository)
	wishlistUseCase := usecase.NewWishlistUseCase(c.DB, c.Log, c.Redis, wishlistRepository)

	// handler
	userHandler := handler.NewUserHandler(c.Log, c.Validator, c.FacebookOAuth, c.GoogleOAuth, userUseCase)
	productHandler := handler.NewProductHandler(c.Log, c.Validator, productUseCase)
	cartHandler := handler.NewCartHandler(c.Log, c.Validator, cartUseCase)
	wishlistHandler := handler.NewWishHandler(c.Log, c.Validator, wishlistUseCase)

	// middleware
	authUserMiddleware := middleware.Auth(c.Log)
	optionalAuthMiddleware := middleware.OptionalAuth()
	corsMiddleware := middleware.Cors()
	httpMiddleware := middleware.HTTP()

	server := &http.Config{
		App:                    c.App,
		UserHandler:            userHandler,
		ProductHandler:         productHandler,
		CartHandler:            cartHandler,
		WishlistHandler:        wishlistHandler,
		AuthUserMiddleware:     authUserMiddleware,
		OptionalAuthMiddleware: optionalAuthMiddleware,
		CorsMiddleware:         corsMiddleware,
		HTTPMiddleware:         httpMiddleware,
	}

	server.Start()
}
