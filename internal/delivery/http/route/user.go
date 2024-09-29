package route

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/interfaces"
	"github.com/gofiber/fiber/v3"
)

func NewUserRoutes(route fiber.Router,
	userHandler interfaces.UserHandler,
	productHandler interfaces.ProductHandler,
	cartHandler interfaces.CartHandler,
	wishlistHandler interfaces.WishlistHandler,
	propertyHandler interfaces.PropertyHandler,
	authMiddleware fiber.Handler,
	optionalAuthMiddleware fiber.Handler) {
	route.Get("/hello", func(ctx fiber.Ctx) error {
		return ctx.SendString("hello user!")
	})

	// @TODO: add (get) renew_access_token and (delete) logout
	auth := route.Group("/auth")
	auth.Get("/facebook", userHandler.URLOAuthFacebook)
	auth.Get("/google", userHandler.URLOAuthGoogle)
	auth.Get("/facebook/callback", userHandler.FacebookOAuthCallback)
	auth.Get("/google/callback", userHandler.GoogleOAuthCallback)
	auth.Post("/oauth/register", userHandler.RegisterWithOAuth)
	auth.Post("/login", userHandler.Login)
	auth.Post("/pre-register", userHandler.PreRegister)
	auth.Post("/register", userHandler.RegisterWithEmailPassword)
	auth.Post("/register/send", userHandler.SendVerificationCodeForRegister)
	auth.Post("/register/complete", userHandler.RegisterComplete)
	auth.Get("/renew-access-token", userHandler.RenewAccessToken)
	auth.Delete("/logout", userHandler.Logout)

	product := route.Group("/product")
	product.Get("/cart", cartHandler.GetMyCart, authMiddleware)
	product.Put("/cart", cartHandler.ManageCart, authMiddleware)
	product.Get("/wishlist", wishlistHandler.GetMyWishlist, authMiddleware)
	product.Put("/wishlist", wishlistHandler.ManageWishlists, authMiddleware)
	product.Get("/:categoryName?", productHandler.GetProducts, optionalAuthMiddleware)
	product.Get("/:productID/details", productHandler.GetProductDetails, optionalAuthMiddleware)
	product.Get("/:productID/reviews", productHandler.GetProductReviews, optionalAuthMiddleware)

	property := route.Group("/property")
	property.Get("/wishlist", propertyHandler.GetMyWishlistProperties, authMiddleware)
	property.Put("/wishlist", propertyHandler.ManageWishlistsProperties, authMiddleware)
	property.Get("/:categoryName?", propertyHandler.GetProperties, optionalAuthMiddleware)
	property.Get("/:propertyID/details", propertyHandler.GetPropertyDetails, optionalAuthMiddleware)
	property.Get("/:propertyID/reviews", propertyHandler.GetPropertyReviews, optionalAuthMiddleware)
	property.Get("/:propertyID/discuss", propertyHandler.GetPropertyDiscuss, optionalAuthMiddleware)
	property.Post("/:propertyID/discuss", propertyHandler.AddPropertyDiscuss, authMiddleware)
}
