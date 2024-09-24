package repository

import (
	"github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
)

type RatingMediaRepository struct {
	DB *sqlx.DB
}

func NewRatingMediaRepository(DB *sqlx.DB) interfaces.RatingMediaRepository {
	return &RatingMediaRepository{
		DB: DB,
	}
}
