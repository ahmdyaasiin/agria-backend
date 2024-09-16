package interfaces

import (
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	//
	Create(tx *sqlx.Tx, user *domain.User) error
	Read(tx *sqlx.Tx, key string, user *domain.User) error
	Update(tx *sqlx.Tx, user *domain.User) error
	Delete(tx *sqlx.Tx, user *domain.User) error

	CheckUserExists(tx *sqlx.Tx, user *domain.User) error
}
