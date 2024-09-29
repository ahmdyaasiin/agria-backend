package http

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/interfaces"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/route"
	"github.com/gofiber/fiber/v3"
)

type Config struct {
	App                    *fiber.App
	AuthUserMiddleware     fiber.Handler
	OptionalAuthMiddleware fiber.Handler
	CorsMiddleware         fiber.Handler
	HTTPMiddleware         fiber.Handler
	UserHandler            interfaces.UserHandler
	ProductHandler         interfaces.ProductHandler
	CartHandler            interfaces.CartHandler
	WishlistHandler        interfaces.WishlistHandler
	PropertyHandler        interfaces.PropertyHandler
	MenuHandler            interfaces.MenuHandler
}

func (c *Config) Start() {
	c.App.Use(c.CorsMiddleware)
	c.App.Use(c.HTTPMiddleware)

	v1 := c.App.Group("/v1")

	route.NewPingRoute(v1)
	route.NewUserRoutes(v1, c.UserHandler, c.ProductHandler, c.CartHandler, c.WishlistHandler, c.PropertyHandler, c.MenuHandler, c.AuthUserMiddleware, c.OptionalAuthMiddleware)
	route.NewAdminRoutes(v1.Group("/admin"))
}
