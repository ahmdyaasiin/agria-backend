package interfaces

import (
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	//
	Create(tx *sqlx.Tx, user *domain.User) error
}
