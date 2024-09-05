package handler

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/interfaces"
	usecaseInterface "github.com/ahmdyaasiin/agria-backend/internal/usecase/interfaces"
)

type UserHandler struct {
	userUserCase usecaseInterface.UserUseCase
}

func NewUserHandler(userUC usecaseInterface.UserUseCase) interfaces.UserHandler {
	return &UserHandler{userUserCase: userUC}
}
