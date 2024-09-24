package interfaces

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/jmoiron/sqlx"
)

type RatingRepository interface {
	//
	GetProductReviews(tx *sqlx.Tx, productID, userID, sortBy string, page int, reviews *[]response.Review) error
	CountRating(tx *sqlx.Tx, productID string, countRatings *int64) error
	RatingBreakdown(tx *sqlx.Tx, productID string, ratingBreakdown *[]response.RatingBreakdown) error
}
