package repository

import (
	"fmt"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/query"
	"github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	DB *sqlx.DB
}

func NewProductRepository(DB *sqlx.DB) interfaces.ProductRepository {
	return &ProductRepository{
		DB: DB,
	}
}

func (r *ProductRepository) Create(tx *sqlx.Tx, product *domain.Product) error {
	_, err := tx.NamedExec(query.CreateQueryBuilder(product), product)
	return err
}

func (r *ProductRepository) Read(tx *sqlx.Tx, key string, product *domain.Product) error {
	q := query.ReadQueryBuilder(product, key)

	value, err := query.GetValueByKey(product, key)
	if err != nil {
		return err
	}

	param := map[string]any{
		key: value,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Get(product, param)
	if err != nil {
		return err
	}

	return err
}

func (r *ProductRepository) Update(tx *sqlx.Tx, product *domain.Product) error {
	_, err := tx.NamedExec(query.UpdateQueryBuilder(product), product)
	return err
}

func (r *ProductRepository) Delete(tx *sqlx.Tx, product *domain.Product) error {
	_, err := tx.NamedExec(query.DeleteQueryBuilder(product), product)
	return err
}

func (r *ProductRepository) GetAllProductsWithoutPromo(tx *sqlx.Tx, categoryName, userID, sortBy, notIN string, page, limit int, product *[]response.GetProduct) error {
	q := QueryGetAllProductsWithoutCondition1
	param := map[string]any{}

	if userID != "" {
		q += QueryAdditionalForGetAllProducts
		param["user_id"] = userID
	}

	q += QueryGetAllProductsWithoutCondition2

	if categoryName != "" {
		q += " AND c.name = :category_name"
		param["category_name"] = categoryName
	}

	q += fmt.Sprintf(" AND p.id NOT IN %s GROUP BY p.id", notIN)

	switch sortBy {
	case "high_rating":
		q += ", ratings ORDER BY ratings DESC"
	case "high_price":
		q += ", p.price ORDER BY p.price DESC"
	case "low_price":
		q += ", p.price ORDER BY p.price"
	default:
		q += ", p.created_at ORDER BY p.created_at DESC"
	}

	q += fmt.Sprintf(" LIMIT %d, %d", (page-1)*5, limit)

	fmt.Println(q)

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Select(product, param)
	if err != nil {
		return err
	}

	return err
}

func (r *ProductRepository) GetAllProductsWithPromo(tx *sqlx.Tx, categoryName, userID, sortBy, notIN string, page int, product *[]response.GetProduct) error {
	q := QueryGetAllProductsWithoutCondition1
	param := map[string]any{}

	if userID != "" {
		q += QueryAdditionalForGetAllProducts
		param["user_id"] = userID
	}

	q += QueryGetAllProductsWithoutCondition2

	if categoryName != "" {
		q += " AND c.name = :category_name"
		param["category_name"] = categoryName
	}

	q += fmt.Sprintf(" AND p.id IN %s GROUP BY p.id", notIN)

	switch sortBy {
	case "high_rating":
		q += ", ratings ORDER BY ratings DESC"
	case "high_price":
		q += ", p.price ORDER BY p.price DESC"
	case "low_price":
		q += ", p.price ORDER BY p.price"
	default:
		q += ", p.created_at ORDER BY p.created_at DESC"
	}

	if page != 0 {
		q += fmt.Sprintf(" LIMIT %d, 24", (page-1)*5)
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Select(product, param)
	if err != nil {
		return err
	}

	return err
}

func (r *ProductRepository) GetDetailsProduct(tx *sqlx.Tx, productID, userID string, product *response.GetProductDetails) error {
	q := QueryGetDetailsProduct1

	param := map[string]any{
		"product_id": productID,
	}
	if userID != "" {
		q += QueryAdditionalForGetDetailsProduct
		param["user_id"] = userID
	}

	q += QueryGetDetailsProduct2

	fmt.Println(q)

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Get(product, param)
	if err != nil {
		return err
	}

	return err
}
