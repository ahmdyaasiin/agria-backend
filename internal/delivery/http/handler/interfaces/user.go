package interfaces

import "github.com/gofiber/fiber/v3"

type UserHandler interface {
	//
	URLOAuthFacebook(ctx fiber.Ctx) error
	URLOAuthGoogle(ctx fiber.Ctx) error

	FacebookOAuthCallback(ctx fiber.Ctx) error
	GoogleOAuthCallback(ctx fiber.Ctx) error

	RegisterWithOAuth(ctx fiber.Ctx) error

	PreRegister(ctx fiber.Ctx) error
	RegisterWithEmailPassword(ctx fiber.Ctx) error
	SendVerificationCodeForRegister(ctx fiber.Ctx) error
	RegisterComplete(ctx fiber.Ctx) error

	Login(ctx fiber.Ctx) error
}
