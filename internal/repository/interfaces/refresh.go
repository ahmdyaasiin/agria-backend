package interfaces

import (
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type RefreshRepository interface {
	//
	Create(tx *sqlx.Tx, refresh *domain.Refresh) error
	Read(tx *sqlx.Tx, key string, refresh *domain.Refresh) error
	Update(tx *sqlx.Tx, refresh *domain.Refresh) error
	Delete(tx *sqlx.Tx, refresh *domain.Refresh) error

	Count(tx *sqlx.Tx, key string, total *int, refresh *domain.Refresh) error
	ReadDESC(tx *sqlx.Tx, key string, refresh *domain.Refresh) error
}
