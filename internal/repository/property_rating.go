package repository

import (
	"github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
)

type PropertyRatingRepository struct {
	DB *sqlx.DB
}

func NewPropertyRatingRepository(DB *sqlx.DB) interfaces.PropertyRatingRepository {
	return &PropertyRatingRepository{DB: DB}
}
