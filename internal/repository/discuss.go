package repository

import (
	"github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
)

type DiscussRepository struct {
	DB *sqlx.DB
}

func NewDiscussRepository(DB *sqlx.DB) interfaces.DiscussRepository {
	return &DiscussRepository{DB: DB}
}
