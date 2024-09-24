package interfaces

import (
	"context"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
)

type ProductUseCase interface {
	//
	GetProducts(ctx context.Context, userID, categoryName, sortBy string, page int) (*response.GetProductWithPagination, error)
	GetProductDetails(ctx context.Context, userID, productID string) (*response.GetProductDetails, error)

	GetProductReviews(ctx context.Context, userID, productID, sortBy string, page int) (*response.ReviewDetails, error)
}
