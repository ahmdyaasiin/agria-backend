package interfaces

import "github.com/gofiber/fiber/v3"

type CartHandler interface {
	//
	GetMyCart(ctx fiber.Ctx) error
	ManageCart(ctx fiber.Ctx) error
}
