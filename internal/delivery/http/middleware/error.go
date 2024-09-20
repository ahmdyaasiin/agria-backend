package middleware

import "github.com/gofiber/fiber/v3"

var (
	ErrNeedBearerToken = fiber.NewError(fiber.StatusBadRequest, "value of authorization header should 'Bearer your_token'")
	ErrInvalidToken    = fiber.NewError(fiber.StatusUnauthorized, "invalid token")
)
