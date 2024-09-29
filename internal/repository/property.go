package repository

import (
	"fmt"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
)

type PropertyRepository struct {
	DB *sqlx.DB
}

func NewPropertyRepository(DB *sqlx.DB) interfaces.PropertyRepository {
	return &PropertyRepository{DB: DB}
}

func (r *PropertyRepository) GetAllPropertiesWithoutPromo(tx *sqlx.Tx, categoryName, userID, sortBy, notIN, province string, page int, properties *[]response.GetProperties) error {
	q := QueryGetAllPropertiesWithoutCondition1
	param := map[string]any{}

	if userID != "" {
		q += QueryAdditionalForGetAllProperties
		param["user_id"] = userID
	}

	q += QueryGetAllPropertiesWithoutCondition2

	if categoryName != "" {
		q += " AND c.name = :category_name"
		param["category_name"] = categoryName
	}

	if province != "all" {
		q += " AND p.state = :state"
		param["state"] = province
	}

	if notIN != "()" {
		q += fmt.Sprintf(" AND p.id NOT IN %s", notIN)
	}

	q += " GROUP BY p.id"

	switch sortBy {
	case "high_rating":
		q += ", ratings ORDER BY ratings DESC"
	case "high_price":
		q += ", p.price ORDER BY p.price DESC"
	case "low_price":
		q += ", p.price ORDER BY p.price"
	default:
		q += ", p.created_at ORDER BY p.created_at DESC"
	}

	if page != 0 {
		q += fmt.Sprintf(" LIMIT %d, 5", (page-1)*5)
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Select(properties, param)
	if err != nil {
		return err
	}

	return err
}

func (r *PropertyRepository) GetAllPropertiesWithPromo(tx *sqlx.Tx, categoryName, userID, sortBy, notIN, province string, page int, properties *[]response.GetProperties) error {
	q := QueryGetAllPropertiesWithoutCondition1
	param := map[string]any{}

	if userID != "" {
		q += QueryAdditionalForGetAllProperties
		param["user_id"] = userID
	}

	q += QueryGetAllPropertiesWithoutCondition2

	if categoryName != "" {
		q += " AND c.name = :category_name"
		param["category_name"] = categoryName
	}

	if province != "all" {
		q += " AND p.state = :state"
		param["state"] = province
	}

	if notIN != "()" {
		q += fmt.Sprintf(" AND p.id IN %s", notIN)
	}

	q += " GROUP BY p.id"

	switch sortBy {
	case "high_rating":
		q += ", ratings ORDER BY ratings DESC"
	case "high_price":
		q += ", p.price ORDER BY p.price DESC"
	case "low_price":
		q += ", p.price ORDER BY p.price"
	default:
		q += ", p.created_at ORDER BY p.created_at DESC"
	}

	if page != 0 {
		q += fmt.Sprintf(" LIMIT %d, 5", (page-1)*5)
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Select(properties, param)
	if err != nil {
		return err
	}

	return err
}

func (r *PropertyRepository) GetPropertyDetails(tx *sqlx.Tx, propertyID, userID string, property *response.GetPropertyDetails) error {
	q := QueryGetPropertyDetails

	param := map[string]any{
		"property_id": propertyID,
		"user_id":     userID,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Get(property, param)
	if err != nil {
		return err
	}

	return err
}

func (r *PropertyRepository) GetPropertyHighlights(tx *sqlx.Tx, propertyID string, highlights *[]response.PropertyHighlights) error {
	q := QueryGetPropertyHighlights

	param := map[string]any{
		"property_id": propertyID,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Select(highlights, param)
	if err != nil {
		return err
	}

	return err
}

func (r *PropertyRepository) GetPropertyMedia(tx *sqlx.Tx, propertyID string, highlights *[]string) error {
	q := QueryGetPropertyMedia

	param := map[string]any{
		"property_id": propertyID,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Select(highlights, param)
	if err != nil {
		return err
	}

	return err
}

func (r *PropertyRepository) GetPropertyDiscuss(tx *sqlx.Tx, propertyID string, discuss *[]response.PropertyDiscuss) error {
	q := QueryGetPropertyDiscuss

	param := map[string]any{
		"property_id": propertyID,
	}

	fmt.Println(q)

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Select(discuss, param)
	if err != nil {
		return err
	}

	return err
}

func (r *PropertyRepository) GetPropertyRatings(tx *sqlx.Tx, propertyID, userID string, ratings *[]response.RatingProperty) error {
	q := QueryGetPropertyRatings

	param := map[string]any{
		"property_id": propertyID,
		"user_id":     userID,
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	err = stmt.Select(ratings, param)
	if err != nil {
		return err
	}

	return err
}

func (r *PropertyRepository) GetState(tx *sqlx.Tx, provinces *[]string) error {
	q := QueryGetProvinces

	err := tx.Select(provinces, q)
	if err != nil {
		return err
	}

	return err
}

func (r *PropertyRepository) RatingBreakdown(tx *sqlx.Tx, productID string, ratingBreakdown *[]response.RatingBreakdown) error {
	q := QueryRatingPropertyBreakdown

	param := map[string]any{
		"property_id": productID,
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
