package repository

import (
	"fmt"
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

func (r *EducationRepository) MainArticle(tx *sqlx.Tx, id string, education *response.EducationCard) error {
	q := QueryGetMainArticle

	param := map[string]any{
		"id": id,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Get(education, param)
	if err != nil {
		return err
	}

	return err
}

func (r *EducationRepository) MustRead(tx *sqlx.Tx, ids string, educations *[]response.EducationCard) error {
	q := QueryGetMustRead
	q = fmt.Sprintf(q, ids)

	fmt.Println(q)

	err := tx.Select(educations, q)
	if err != nil {
		return err
	}

	return err
}

func (r *EducationRepository) Latest(tx *sqlx.Tx, ids string, educations *[]response.EducationCard) error {
	q := QueryGetLatest
	q = fmt.Sprintf(q, ids)

	fmt.Println(q)

	err := tx.Select(educations, q)
	if err != nil {
		return err
	}

	return err
}

func (r *EducationRepository) ExceptionWithRandom(tx *sqlx.Tx, ids string, educations *[]response.EducationCard) error {
	q := QueryGetDiscoverMore
	q = fmt.Sprintf(q, ids)

	fmt.Println(q)

	err := tx.Select(educations, q)
	if err != nil {
		return err
	}

	return err
}

func (r *EducationRepository) EducationDetails(tx *sqlx.Tx, id string, education *response.EducationData) error {
	q := QueryGetEducationDetails

	param := map[string]any{
		"id": id,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Get(education, param)
	if err != nil {
		return err
	}

	return err
}
