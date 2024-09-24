package repository

import (
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/query"
	"github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
)

type AddressRepository struct {
	DB *sqlx.DB
}

func NewAddressRepository(DB *sqlx.DB) interfaces.AddressRepository {
	return &AddressRepository{DB: DB}
}

func (r *AddressRepository) Create(tx *sqlx.Tx, address *domain.Address) error {
	_, err := tx.NamedExec(query.CreateQueryBuilder(address), address)
	return err
}

func (r *AddressRepository) Read(tx *sqlx.Tx, key string, address *domain.Address) error {
	q := query.ReadQueryBuilder(address, key)
	q += " AND is_primary = 1"

	value, err := query.GetValueByKey(address, key)
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

	err = stmt.Get(address, param)
	if err != nil {
		return err
	}

	return err
}

func (r *AddressRepository) Update(tx *sqlx.Tx, address *domain.Address) error {
	_, err := tx.NamedExec(query.UpdateQueryBuilder(address), address)
	return err
}

func (r *AddressRepository) Delete(tx *sqlx.Tx, address *domain.Address) error {
	_, err := tx.NamedExec(query.DeleteQueryBuilder(address), address)
	return err
}
