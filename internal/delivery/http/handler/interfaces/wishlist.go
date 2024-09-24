package interfaces

import "github.com/gofiber/fiber/v3"

type WishlistHandler interface {
	//
	GetMyWishlist(ctx fiber.Ctx) error
	ManageWishlists(ctx fiber.Ctx) error
}
