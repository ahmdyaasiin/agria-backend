package interfaces

import (
	"context"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
)

type MenuUseCase interface {
	//
	Homepage(ctx context.Context, auth string) (*response.Homepage, error)
	Market(ctx context.Context, userID string) (*response.Market, error)
}
