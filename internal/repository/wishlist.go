package repository

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/query"
	"github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
)

type WishlistRepository struct {
	DB *sqlx.DB
}

func NewWishlistRepository(DB *sqlx.DB) interfaces.WishlistRepository {
	return &WishlistRepository{DB: DB}
}

func (r *WishlistRepository) GetMyWishlists(tx *sqlx.Tx, userID string, wishlists *[]response.MyWishlist) error {
	q := QueryGetMyWishlist

	param := map[string]any{
		"user_id": userID,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Select(wishlists, param)
	if err != nil {
		return err
	}

	return err
}

func (r *WishlistRepository) GetSpecificProduct(tx *sqlx.Tx, userID, productID string, wishlist *domain.Wishlist) error {
	q := QueryGetSpecificProductInWishlist

	param := map[string]any{
		"user_id":    userID,
		"product_id": productID,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Get(wishlist, param)
	if err != nil {
		return err
	}

	return err
}

func (r *WishlistRepository) Create(tx *sqlx.Tx, wishlist *domain.Wishlist) error {
	_, err := tx.NamedExec(query.CreateQueryBuilder(wishlist), wishlist)
	return err
}

func (r *WishlistRepository) Delete(tx *sqlx.Tx, wishlist *domain.Wishlist) error {
	_, err := tx.NamedExec(query.DeleteQueryBuilder(wishlist), wishlist)
	return err
}
