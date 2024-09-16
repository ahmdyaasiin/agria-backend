package repository

import (
	"database/sql"
	"errors"
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/query"
	"github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(DB *sqlx.DB) interfaces.UserRepository {
	return &UserRepository{DB: DB}
}

func (r *UserRepository) Create(tx *sqlx.Tx, user *domain.User) error {
	_, err := tx.NamedExec(query.CreateQueryBuilder(user), user)
	return err
}

func (r *UserRepository) Read(tx *sqlx.Tx, key string, user *domain.User) error {
	q := query.ReadQueryBuilder(user, key)

	value, err := query.GetValueByKey(user, key)
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

	err = stmt.Get(user, param)
	if err != nil {
		return err
	}

	return err
}

func (r *UserRepository) Update(tx *sqlx.Tx, user *domain.User) error {
	_, err := tx.NamedExec(query.UpdateQueryBuilder(user), user)
	return err
}

func (r *UserRepository) Delete(tx *sqlx.Tx, user *domain.User) error {
	_, err := tx.NamedExec(query.DeleteQueryBuilder(user), user)
	return err
}

func (u *UserRepository) CheckUserExists(tx *sqlx.Tx, user *domain.User) error {
	err := u.Read(tx, "email", user)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if user.ID != "" {
		return errors.New("user already registered with this email")
	}

	err = u.Read(tx, "username", user)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if user.ID != "" {
		return errors.New("user already registered with this username")
	}

	err = u.Read(tx, "phone_number", user)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if user.ID != "" {
		return errors.New("user already registered with this phone number")
	}

	return nil
}
