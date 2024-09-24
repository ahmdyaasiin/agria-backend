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
	"strconv"
)

type ProductHandler struct {
	Log            *logrus.Logger
	Validator      *validator.Validate
	ProductUseCase usecaseInterface.ProductUseCase
}

func NewProductHandler(log *logrus.Logger, validator *validator.Validate, productUC usecaseInterface.ProductUseCase) interfaces.ProductHandler {

	return &ProductHandler{
		Log:            log,
		Validator:      validator,
		ProductUseCase: productUC,
	}
}

func (h *ProductHandler) GetProducts(ctx fiber.Ctx) error {
	auth := middleware.GetUserID(ctx)
	category := ctx.Params("categoryName")
	sortBy := ctx.Query("sortBy", "newest")
	pageString := ctx.Query("page", "0")

	page, err := strconv.Atoi(pageString)
	if err != nil {
		return err
	}

	if page < 0 || (sortBy != "newest" && sortBy != "high_rating" && sortBy != "high_price" && sortBy != "low_price") {
		return fiber.NewError(fiber.StatusBadRequest, "query string is invalid")
	}

	res, err := h.ProductUseCase.GetProducts(ctx.Context(), auth, category, sortBy, page)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "Get product successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}

func (h *ProductHandler) GetProductDetails(ctx fiber.Ctx) error {
	auth := middleware.GetUserID(ctx)
	productID := ctx.Params("productID")

	res, err := h.ProductUseCase.GetProductDetails(ctx.Context(), auth, productID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "Get product details successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}

func (h *ProductHandler) GetProductReviews(ctx fiber.Ctx) error {
	auth := middleware.GetUserID(ctx)
	productID := ctx.Params("productID")
	sortBy := ctx.Query("sortBy", "newest")
	pageString := ctx.Query("page", "0")

	page, err := strconv.Atoi(pageString)
	if err != nil {
		return err
	}

	if page < 0 || (sortBy != "newest" && sortBy != "helpful" && sortBy != "high_rating" && sortBy != "low_rating") {
		return fiber.NewError(fiber.StatusBadRequest, "query string is invalid")
	}

	res, err := h.ProductUseCase.GetProductReviews(ctx.Context(), auth, productID, sortBy, page)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "Get product reviews successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}
