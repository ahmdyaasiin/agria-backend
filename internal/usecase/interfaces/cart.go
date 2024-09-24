package interfaces

import (
	"context"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/request"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
)

type CartUseCase interface {
	//
	GetMyCart(ctx context.Context, auth string) (*response.MyCart, error)
	ManageCart(ctx context.Context, userID string, req *request.ManageCart) (*response.ManageCart, error)
}
