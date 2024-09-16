package handler

import "github.com/gofiber/fiber/v3"

// general error

// user handler error
var (
	ErrFRExchangeOAuthCode = "failed+to+extract+code+parameter"
	ErrFRFetchOAuthProfile = "failed+to+fetch+profile"
)

var (
	ErrBindRequest     = fiber.NewError(fiber.StatusBadRequest, "failed to bind request")
	ErrValidateRequest = fiber.NewError(fiber.StatusBadRequest, "failed to validate request")
)
