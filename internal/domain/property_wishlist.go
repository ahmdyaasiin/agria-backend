package domain

type PropertyWishlist struct {
	//
	ID         string `json:"id" db:"id"`
	CreatedAt  int64  `json:"created_at" db:"created_at"`
	UserID     string `json:"user_id" db:"user_id"`
	PropertyID string `json:"property_id" db:"property_id"`
}

func (e PropertyWishlist) TableName() string {
	return "property_wishlist"
}
