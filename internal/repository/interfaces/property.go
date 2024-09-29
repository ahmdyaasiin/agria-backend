package interfaces

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/jmoiron/sqlx"
)

type PropertyRepository interface {
	//
	GetAllPropertiesWithoutPromo(tx *sqlx.Tx, categoryName, userID, sortBy, notIN, province string, page int, properties *[]response.GetProperties) error
	GetPropertyDetails(tx *sqlx.Tx, propertyID, userID string, property *response.GetPropertyDetails) error

	GetPropertyHighlights(tx *sqlx.Tx, propertyID string, highlights *[]response.PropertyHighlights) error
	GetPropertyMedia(tx *sqlx.Tx, propertyID string, highlights *[]string) error

	GetPropertyDiscuss(tx *sqlx.Tx, propertyID string, discuss *[]response.PropertyDiscuss) error

	GetPropertyRatings(tx *sqlx.Tx, propertyID, userID string, ratings *[]response.RatingProperty) error
}
