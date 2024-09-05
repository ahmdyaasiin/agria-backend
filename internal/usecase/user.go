package usecase

import (
	repositoryInterface "github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/ahmdyaasiin/agria-backend/internal/usecase/interfaces"
	"github.com/jmoiron/sqlx"
)

type userUseCase struct {
	DB          *sqlx.DB
	userUseCase repositoryInterface.UserRepository
}

func NewUserUseCase(DB *sqlx.DB, userUC repositoryInterface.UserRepository) interfaces.UserUseCase {
	return &userUseCase{DB: DB, userUseCase: userUC}
}
