package interfaces

import (
	"github.com/jmoiron/sqlx"
)

type ProductMediaRepository interface {
	//
	GetProductMedia(tx *sqlx.Tx, productID string, productMedia *[]string) error
}
