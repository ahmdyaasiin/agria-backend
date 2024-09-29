package interfaces

import (
	"context"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/request"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
)

type PropertyUseCase interface {
	//
	GetAllWishlistsProperties(ctx context.Context, userID string) (*response.PropertiesWishlist, error)
	ManageWishlistProperties(ctx context.Context, userID string, req *request.ManageWishlistProperties) (*response.ManageWishlistProperties, error)

	GetProperties(ctx context.Context, userID, categoryName, sortBy, province string, page int) (*response.GetPropertiesWithPagination, error)
	GetPropertyDetails(ctx context.Context, userID, propertyID string) (*response.GetPropertyDetails, error)
}
