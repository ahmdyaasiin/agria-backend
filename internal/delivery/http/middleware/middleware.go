package middleware

import (
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/sirupsen/logrus"
	"strings"
)

func Auth(log *logrus.Logger) fiber.Handler {
	return func(ctx fiber.Ctx) error {

		authorizationHeader := ctx.GetReqHeaders()["Authorization"]

		var bearer string
		if len(authorizationHeader) > 0 {
			bearer = authorizationHeader[0]
		} else {
			log.Warnf("no authorization header is sent: %+v\n", bearer)
			return ErrNeedBearerToken
		}

		if strings.HasPrefix(bearer, "Bearer ") {
			bearer = strings.Split(bearer, " ")[1]
		} else {
			log.Warnf("no bearer prefix at the header: %+v\n", bearer)
			return ErrNeedBearerToken
		}

		userID, err := jwt.ValidateToken(bearer, false)
		if err != nil {
			log.Warnf("failed to validate token: %+v\n", err)
			return ErrInvalidToken
		}

		ctx.Locals("user_id", userID)
		return ctx.Next()
	}
}

func OptionalAuth() fiber.Handler {
	return func(ctx fiber.Ctx) error {

		authorizationHeader := ctx.GetReqHeaders()["Authorization"]

		var bearer string
		if len(authorizationHeader) > 0 {
			bearer = authorizationHeader[0]

			if strings.HasPrefix(bearer, "Bearer ") {
				bearer = strings.Split(bearer, " ")[1]
				userID, _ := jwt.ValidateToken(bearer, false)

				ctx.Locals("user_id", userID)
			}
		}

		return ctx.Next()
	}
}

func Cors() fiber.Handler {
	return cors.New()
}

func HTTP() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${method} ${path} ${status} - ${time} - ${latency}\n",
		TimeFormat: "15:04:05 Jan 2 2006",
		TimeZone:   "Local",
	})
}

func GetUserID(ctx fiber.Ctx) string {
	if userID, ok := ctx.Locals("user_id").(string); ok {
		return userID
	}

	return ""
}
