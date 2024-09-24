package domain

type Wishlist struct {
	//
	ID        string `json:"id" db:"id"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
	UserID    string `json:"user_id" db:"user_id"`
	ProductID string `json:"product_id" db:"product_id"`
}

func (e Wishlist) TableName() string {
	return "wishlists"
}
