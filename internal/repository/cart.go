package repository

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/query"
	"github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
)

type CartRepository struct {
	DB *sqlx.DB
}

func NewCartRepository(DB *sqlx.DB) interfaces.CartRepository {
	return &CartRepository{DB: DB}
}

func (r *CartRepository) Create(tx *sqlx.Tx, cart *domain.Cart) error {
	_, err := tx.NamedExec(query.CreateQueryBuilder(cart), cart)
	return err
}

func (r *CartRepository) Update(tx *sqlx.Tx, cart *domain.Cart) error {
	_, err := tx.NamedExec(query.UpdateQueryBuilder(cart), cart)
	return err
}

func (r *CartRepository) Delete(tx *sqlx.Tx, cart *domain.Cart) error {
	_, err := tx.NamedExec(query.DeleteQueryBuilder(cart), cart)
	return err
}

func (r *CartRepository) GetMyCartAvailable(tx *sqlx.Tx, userID string, carts *[]response.CartProducts) error {
	q := QueryGetMyCartAvailable

	param := map[string]any{
		"user_id": userID,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Select(carts, param)
	if err != nil {
		return err
	}

	return err
}

func (r *CartRepository) GetMyCartUnavailable(tx *sqlx.Tx, userID string, carts *[]response.CartProducts) error {
	q := QueryGetMyCartUnavailable

	param := map[string]any{
		"user_id": userID,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Select(carts, param)
	if err != nil {
		return err
	}

	return err
}

func (r *CartRepository) GetMyCart(tx *sqlx.Tx, userID, productID string, cart *domain.Cart) error {
	q := QueryGetSpecificProductInCart

	param := map[string]any{
		"product_id": productID,
		"user_id":    userID,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Get(cart, param)
	if err != nil {
		return err
	}

	return err
}

func (r *CartRepository) CountCart(tx *sqlx.Tx, userID string, total *int) error {
	q := QueryGetCountCart

	param := map[string]any{
		"user_id": userID,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Get(total, param)
	if err != nil {
		return err
	}

	return err
}
