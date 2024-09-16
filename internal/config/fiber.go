package config

import (
	"errors"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"net/http"
	"os"
)

func NewFiber() *fiber.App {

	app := fiber.New(fiber.Config{
		AppName:      os.Getenv("APP_NAME"),
		ErrorHandler: NewErrorHandler(),
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8000"},
		AllowCredentials: true,
	}))

	app.Use(logger.New(logger.Config{
		Format:     "${method} ${path} ${status} - ${time} - ${latency}\n",
		TimeFormat: "15:04:05 Jan 2 2006",
		TimeZone:   "Local",
	}))

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx fiber.Ctx, err error) error {

		code := fiber.StatusInternalServerError
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorList := validation.GetError(err, ve)
			if err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(response.Final{
					Message: "validation error",
					Data:    nil,
					Errors:  errorList,
					Status: response.Status{
						Code:    fiber.StatusBadRequest,
						Message: http.StatusText(fiber.StatusBadRequest),
					},
				})
			}
		}

		return ctx.Status(code).JSON(response.Final{
			Message: err.Error(),
			Data:    nil,
			Errors:  nil,
			Status: response.Status{
				Code:    code,
				Message: http.StatusText(code),
			},
		})
	}
}
