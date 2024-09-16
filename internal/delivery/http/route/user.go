package route

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/interfaces"
	"github.com/gofiber/fiber/v3"
)

func NewUserRoutes(route fiber.Router, userHandler interfaces.UserHandler) {
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

}
