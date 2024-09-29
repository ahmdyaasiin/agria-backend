package interfaces

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/jmoiron/sqlx"
)

type EducationRepository interface {
	//
	GetAllEducation(tx *sqlx.Tx, educations *[]response.EducationsHomePage) error
}
