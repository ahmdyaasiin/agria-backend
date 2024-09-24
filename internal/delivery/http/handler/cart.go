package handler

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/interfaces"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/request"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/middleware"
	usecaseInterface "github.com/ahmdyaasiin/agria-backend/internal/usecase/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CartHandler struct {
	Log         *logrus.Logger
	Validator   *validator.Validate
	CartUseCase usecaseInterface.CartUseCase
}

func NewCartHandler(log *logrus.Logger, validator *validator.Validate, cartUseCase usecaseInterface.CartUseCase) interfaces.CartHandler {

	return &CartHandler{
		Log:         log,
		Validator:   validator,
		CartUseCase: cartUseCase,
	}
}

func (h *CartHandler) GetMyCart(ctx fiber.Ctx) error {
	auth := middleware.GetUserID(ctx)

	res, err := h.CartUseCase.GetMyCart(ctx.Context(), auth)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "Get cart successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}

func (h *CartHandler) ManageCart(ctx fiber.Ctx) error {
	auth := middleware.GetUserID(ctx)

	req := new(request.ManageCart)

	if err := ctx.Bind().JSON(req); err != nil {
		h.Log.Warnf("failed to bind request: %+v\n", err)
		return ErrBindRequest
	}

	if err := h.Validator.Struct(req); err != nil {
		h.Log.Warnf("failed to validate request: %+v\n", err)
		return err
	}

	res, err := h.CartUseCase.ManageCart(ctx.Context(), auth, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "Manage cart successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}
