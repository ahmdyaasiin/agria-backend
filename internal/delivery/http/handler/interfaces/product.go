package interfaces

import "github.com/gofiber/fiber/v3"

type ProductHandler interface {
	//
	GetProducts(ctx fiber.Ctx) error
	GetProductDetails(ctx fiber.Ctx) error
	GetProductReviews(ctx fiber.Ctx) error
}
