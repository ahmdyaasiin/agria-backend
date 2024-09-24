package repository

import (
	"github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
)

type ProductMediaRepository struct {
	DB *sqlx.DB
}

func NewProductMediaRepository(DB *sqlx.DB) interfaces.ProductMediaRepository {
	return &ProductMediaRepository{
		DB: DB,
	}
}

func (r *ProductMediaRepository) GetProductMedia(tx *sqlx.Tx, productID string, productMedia *[]string) error {
	q := QueryGetProductMedia

	param := map[string]any{
		"product_id": productID,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Select(productMedia, param)
	if err != nil {
		return err
	}

	return err
}
