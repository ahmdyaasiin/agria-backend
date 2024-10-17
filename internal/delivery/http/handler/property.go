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
	"strconv"
)

type PropertyHandler struct {
	Log             *logrus.Logger
	Validator       *validator.Validate
	PropertyUseCase usecaseInterface.PropertyUseCase
}

func NewPropertyHandler(log *logrus.Logger, validator *validator.Validate, propertyUseCase usecaseInterface.PropertyUseCase) interfaces.PropertyHandler {
	return &PropertyHandler{Log: log, Validator: validator, PropertyUseCase: propertyUseCase}
}

func (h *PropertyHandler) GetMyWishlistProperties(ctx fiber.Ctx) error {
	auth := middleware.GetUserID(ctx)

	res, err := h.PropertyUseCase.GetAllWishlistsProperties(ctx.Context(), auth)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "Get wishlists successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}

func (h *PropertyHandler) ManageWishlistsProperties(ctx fiber.Ctx) error {
	auth := middleware.GetUserID(ctx)

	req := new(request.ManageWishlistProperties)

	if err := ctx.Bind().JSON(req); err != nil {
		h.Log.Warnf("failed to bind request: %+v\n", err)
		return ErrBindRequest
	}

	if err := h.Validator.Struct(req); err != nil {
		h.Log.Warnf("failed to validate request: %+v\n", err)
		return err
	}

	res, err := h.PropertyUseCase.ManageWishlistProperties(ctx.Context(), auth, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "Manage wishlists successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}

func (h *PropertyHandler) GetProperties(ctx fiber.Ctx) error {
	auth := middleware.GetUserID(ctx)
	category := ctx.Params("categoryName")
	sortBy := ctx.Query("sortBy", "newest")
	pageString := ctx.Query("page", "1")
	province := ctx.Query("province", "all")

	page, err := strconv.Atoi(pageString)
	if err != nil {
		return err
	}

	if page < 0 || (sortBy != "newest" && sortBy != "high_rating" && sortBy != "high_price" && sortBy != "low_price") {
		return fiber.NewError(fiber.StatusBadRequest, "query string is invalid")
	}

	res, err := h.PropertyUseCase.GetProperties(ctx.Context(), auth, category, sortBy, province, page)
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

func (h *PropertyHandler) GetPropertyDetails(ctx fiber.Ctx) error {
	auth := middleware.GetUserID(ctx)
	propertyID := ctx.Params("propertyID")

	res, err := h.PropertyUseCase.GetPropertyDetails(ctx.Context(), auth, propertyID)
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

func (h *PropertyHandler) GetPropertyReviews(ctx fiber.Ctx) error {

	return nil
	//auth := middleware.GetUserID(ctx)
	//productID := ctx.Params("productID")
	//sortBy := ctx.Query("sortBy", "newest")
	//pageString := ctx.Query("page", "0")
	//
	//page, err := strconv.Atoi(pageString)
	//if err != nil {
	//	return err
	//}
	//
	//if page < 0 || (sortBy != "newest" && sortBy != "helpful" && sortBy != "high_rating" && sortBy != "low_rating") {
	//	return fiber.NewError(fiber.StatusBadRequest, "query string is invalid")
	//}
	//
	//res, err := h.PropertyUseCase.GetProductReviews(ctx.Context(), auth, productID, sortBy, page)
	//if err != nil {
	//	return err
	//}
	//
	//return ctx.Status(fiber.StatusOK).JSON(response.Final{
	//	Message: "Get product reviews successfully",
	//	Data:    res,
	//	Errors:  nil,
	//	Status: response.Status{
	//		Code:    fiber.StatusOK,
	//		Message: http.StatusText(fiber.StatusOK),
	//	},
	//})
}

func (h *PropertyHandler) GetPropertyDiscuss(ctx fiber.Ctx) error {
	return nil
}

func (h *PropertyHandler) AddPropertyDiscuss(ctx fiber.Ctx) error {
	return nil
}
