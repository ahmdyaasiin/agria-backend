package repository

import (
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(DB *sqlx.DB) interfaces.UserRepository {
	return &userRepository{DB: DB}
}

func (r *userRepository) Create(tx *sqlx.Tx, user *domain.User) error {

	return nil
}
