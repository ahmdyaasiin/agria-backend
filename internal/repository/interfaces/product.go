package interfaces

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	//
	Create(tx *sqlx.Tx, product *domain.Product) error
	Read(tx *sqlx.Tx, key string, product *domain.Product) error
	Update(tx *sqlx.Tx, product *domain.Product) error
	Delete(tx *sqlx.Tx, product *domain.Product) error

	GetAllProductsWithoutPromo(tx *sqlx.Tx, categoryName, userID, sortBy, notIN string, page int, product *[]response.GetProduct) error
	GetDetailsProduct(tx *sqlx.Tx, productID, userID string, product *response.GetProductDetails) error

	GetAllProductsWithPromo(tx *sqlx.Tx, categoryName, userID, sortBy, notIN string, page int, product *[]response.GetProduct) error
}
