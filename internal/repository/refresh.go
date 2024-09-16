package repository

import (
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/query"
	"github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
)

type RefreshRepository struct {
	DB *sqlx.DB
}

func NewRefreshRepository(DB *sqlx.DB) interfaces.RefreshRepository {
	return &RefreshRepository{DB: DB}
}

func (r *RefreshRepository) Create(tx *sqlx.Tx, refresh *domain.Refresh) error {
	_, err := tx.NamedExec(query.CreateQueryBuilder(refresh), refresh)
	return err
}

func (r *RefreshRepository) Read(tx *sqlx.Tx, key string, refresh *domain.Refresh) error {
	q := query.ReadQueryBuilder(refresh, key)

	value, err := query.GetValueByKey(refresh, key)
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

	err = stmt.Get(refresh, param)
	if err != nil {
		return err
	}

	return err
}

func (r *RefreshRepository) Update(tx *sqlx.Tx, refresh *domain.Refresh) error {
	_, err := tx.NamedExec(query.UpdateQueryBuilder(refresh), refresh)
	return err
}

func (r *RefreshRepository) Delete(tx *sqlx.Tx, refresh *domain.Refresh) error {
	_, err := tx.NamedExec(query.DeleteQueryBuilder(refresh), refresh)
	return err
}
