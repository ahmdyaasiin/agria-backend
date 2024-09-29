package interfaces

import "github.com/gofiber/fiber/v3"

type PropertyHandler interface {
	//
	GetMyWishlistProperties(ctx fiber.Ctx) error
	ManageWishlistsProperties(ctx fiber.Ctx) error

	GetProperties(ctx fiber.Ctx) error
	GetPropertyDetails(ctx fiber.Ctx) error
	GetPropertyReviews(ctx fiber.Ctx) error
	GetPropertyDiscuss(ctx fiber.Ctx) error
	AddPropertyDiscuss(ctx fiber.Ctx) error
}
