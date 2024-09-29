package repository

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
)

type EducationRepository struct {
	DB *sqlx.DB
}

func NewEducationRepository(DB *sqlx.DB) interfaces.EducationRepository {
	return &EducationRepository{DB: DB}
}

func (r *EducationRepository) GetAllEducation(tx *sqlx.Tx, educations *[]response.EducationsHomePage) error {
	q := QueryGetEducation

	err := tx.Select(educations, q)
	if err != nil {
		return err
	}

	return err
}
