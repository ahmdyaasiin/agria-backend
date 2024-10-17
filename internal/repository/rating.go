package repository

import (
	"fmt"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
)

type RatingRepository struct {
	DB *sqlx.DB
}

func NewRatingRepository(DB *sqlx.DB) interfaces.RatingRepository {
	return &RatingRepository{
		DB: DB,
	}
}

func (r *RatingRepository) GetProductReviews(tx *sqlx.Tx, productID, userID, sortBy string, page int, reviews *[]response.Review) error {
	q := QueryGetReviewsProduct1

	param := map[string]any{
		"product_id": productID,
	}
	if userID != "" {
		q += QueryAdditionalForGetReviewsProduct
		param["user_id"] = userID
	}

	q += QueryGetReviewsProduct2

	switch sortBy {
	case "helpful":
		q += ", helpful_count ORDER BY helpful_count DESC"
	case "high_rating":
		q += ", r.star ORDER BY r.star DESC"
	case "low_rating":
		q += ", r.star ORDER BY r.star"
	default:
		q += ", r.created_at ORDER BY r.created_at DESC"
	}

	if page != 0 {
		q += fmt.Sprintf(" LIMIT %d, 5", (page-1)*5)
	}

	fmt.Println(q)

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Select(reviews, param)
	if err != nil {
		return err
	}

	return err
}

func (r *RatingRepository) CountRating(tx *sqlx.Tx, productID string, countRatings *int64) error {
	q := QueryCountRatings

	param := map[string]any{
		"product_id": productID,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Get(countRatings, param)
	if err != nil {
		return err
	}

	return err
}

func (r *RatingRepository) RatingBreakdown(tx *sqlx.Tx, productID string, ratingBreakdown *[]response.RatingBreakdown) error {
	q := QueryRatingBreakdown

	param := map[string]any{
		"product_id": productID,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Select(ratingBreakdown, param)
	if err != nil {
		return err
	}

	return err
}
