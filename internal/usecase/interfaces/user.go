package interfaces

import (
	"context"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/request"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
)

type UserUseCase interface {
	//
	FacebookCallBack(ctx context.Context, profile *response.FetchFacebookProfile) (*response.OAuth, error)
	GoogleCallBack(ctx context.Context, profile *response.FetchGoogleProfile) (*response.OAuth, error)

	RegisterWithOAuth(ctx context.Context, req *request.FinishRegisterOAuth) (*response.FinishRegister, error)

	RegisterWithEmailPassword(ctx context.Context, req *request.Register) error
	SendVerificationCodeForRegister(ctx context.Context, req *request.PostRegister) error
	VerifySixCode(ctx context.Context, req *request.FinishRegister) (*response.FinishRegister, error)

	Login(ctx context.Context, req *request.Login) (*response.FinishRegister, error)
}
