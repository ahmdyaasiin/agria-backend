package interfaces

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type CartRepository interface {
	//
	Create(tx *sqlx.Tx, cart *domain.Cart) error
	Update(tx *sqlx.Tx, cart *domain.Cart) error
	Delete(tx *sqlx.Tx, cart *domain.Cart) error

	GetMyCartAvailable(tx *sqlx.Tx, userID string, carts *[]response.CartProducts) error
	GetMyCartUnavailable(tx *sqlx.Tx, userID string, carts *[]response.CartProducts) error

	GetMyCart(tx *sqlx.Tx, userID, productID string, cart *domain.Cart) error

	CountCart(tx *sqlx.Tx, userID string, total *int) error
}
