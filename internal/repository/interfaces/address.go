package interfaces

import (
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type AddressRepository interface {
	//
	Create(tx *sqlx.Tx, address *domain.Address) error
	Read(tx *sqlx.Tx, key string, address *domain.Address) error
	Update(tx *sqlx.Tx, address *domain.Address) error
	Delete(tx *sqlx.Tx, address *domain.Address) error
}
