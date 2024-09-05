package config

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"net/http"
	"os"
)

func NewFiber() *fiber.App {

	app := fiber.New(fiber.Config{
		AppName:      os.Getenv("APP_NAME"),
		ErrorHandler: NewErrorHandler(),
	})

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx fiber.Ctx, err error) error {

		code := fiber.StatusInternalServerError

		e := new(fiber.Error)
		if errors.As(err, &e) {
			code = e.Code
		}

		return ctx.Status(code).SendString(http.StatusText(code))
	}
}
