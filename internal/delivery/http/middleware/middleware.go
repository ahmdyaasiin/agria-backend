package middleware

import (
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"strings"
)

type Middleware struct {
	Log *logrus.Logger
}

func NewMiddleware(log *logrus.Logger) *Middleware {

	return &Middleware{
		Log: log,
	}
}

func (m *Middleware) Auth() fiber.Handler {
	return func(ctx fiber.Ctx) error {

		authorizationHeader := ctx.GetReqHeaders()["Authorization"]

		var bearer string
		if len(authorizationHeader) > 0 {
			bearer = authorizationHeader[0]
		} else {
			return ErrNeedBearerToken
		}

		if strings.HasPrefix(bearer, "Bearer ") {
			bearer = strings.Split(bearer, " ")[1]
		} else {
			return ErrNeedBearerToken
		}

		userID, err := jwt.ValidateToken(bearer, false)
		if err != nil {
			return ErrInvalidToken
		}

		ctx.Locals("user_id", userID)

		return ctx.Next()
	}
}
