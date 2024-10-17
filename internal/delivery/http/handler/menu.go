package handler

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/interfaces"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/middleware"
	usecaseInterface "github.com/ahmdyaasiin/agria-backend/internal/usecase/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"net/http"
)

type MenuHandler struct {
	Log         *logrus.Logger
	Validator   *validator.Validate
	MenuUseCase usecaseInterface.MenuUseCase
}

func NewMenuHandler(log *logrus.Logger, validator *validator.Validate, menuUseCase usecaseInterface.MenuUseCase) interfaces.MenuHandler {

	return &MenuHandler{
		Log:         log,
		Validator:   validator,
		MenuUseCase: menuUseCase,
	}
}

func (h *MenuHandler) GetHomepage(ctx fiber.Ctx) error {
	auth := middleware.GetUserID(ctx)

	res, err := h.MenuUseCase.Homepage(ctx.Context(), auth)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "Get homepage successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}

func (h *MenuHandler) GetMarket(ctx fiber.Ctx) error {
	auth := middleware.GetUserID(ctx)

	res, err := h.MenuUseCase.Market(ctx.Context(), auth)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "Get market successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}

func (h *MenuHandler) GetEducation(ctx fiber.Ctx) error {
	auth := middleware.GetUserID(ctx)

	res, err := h.MenuUseCase.Education(ctx.Context(), auth)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "Get homepage successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}

func (h *MenuHandler) GetEducationDetails(ctx fiber.Ctx) error {
	auth := middleware.GetUserID(ctx)

	educationID := ctx.Params("educationID")
	res, err := h.MenuUseCase.EducationDetails(ctx.Context(), auth, educationID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "Get education details successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}
