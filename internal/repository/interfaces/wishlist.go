package interfaces

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type WishlistRepository interface {
	//
	Create(tx *sqlx.Tx, wishlist *domain.Wishlist) error
	Delete(tx *sqlx.Tx, wishlist *domain.Wishlist) error

	CreateProperty(tx *sqlx.Tx, wishlist *domain.PropertyWishlist) error
	DeleteProperty(tx *sqlx.Tx, wishlist *domain.PropertyWishlist) error

	GetMyWishlists(tx *sqlx.Tx, userID string, wishlists *[]response.MyWishlist) error
	GetSpecificProduct(tx *sqlx.Tx, userID, productID string, wishlist *domain.Wishlist) error

	GetMyWishlistsProperty(tx *sqlx.Tx, userID string, wishlists *[]response.MyWishlistProperties) error
	GetSpecificProperty(tx *sqlx.Tx, userID, productID string, wishlist *domain.PropertyWishlist) error
}
