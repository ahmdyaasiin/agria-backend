package interfaces

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/jmoiron/sqlx"
)

type EducationRepository interface {
	//
	GetAllEducation(tx *sqlx.Tx, educations *[]response.EducationsHomePage) error

	MainArticle(tx *sqlx.Tx, id string, education *response.EducationCard) error
	MustRead(tx *sqlx.Tx, ids string, educations *[]response.EducationCard) error
	Latest(tx *sqlx.Tx, ids string, educations *[]response.EducationCard) error
	ExceptionWithRandom(tx *sqlx.Tx, ids string, educations *[]response.EducationCard) error

	EducationDetails(tx *sqlx.Tx, id string, education *response.EducationData) error
}
