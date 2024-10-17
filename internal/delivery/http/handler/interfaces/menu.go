package interfaces

import "github.com/gofiber/fiber/v3"

type MenuHandler interface {
	//
	GetHomepage(ctx fiber.Ctx) error
	GetMarket(ctx fiber.Ctx) error
	GetEducation(ctx fiber.Ctx) error

	GetEducationDetails(ctx fiber.Ctx) error
}
