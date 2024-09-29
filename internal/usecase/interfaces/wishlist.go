package interfaces

import (
	"context"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/request"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
)

type WishlistUseCase interface {
	//
	GetAllWishlists(ctx context.Context, userID string) (*response.ProductWishlist, error)
	ManageWishlist(ctx context.Context, userID string, req *request.ManageWishlist) (*response.ManageWishlist, error)
}
