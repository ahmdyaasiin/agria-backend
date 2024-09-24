package response

type MyWishlist struct {
	//
	ProductID       string  `json:"product_id" db:"product_id"`
	ID              string  `json:"id" db:"id"`
	Name            string  `json:"name" db:"name"`
	Price           int64   `json:"price" db:"price"`
	DiscountPrice   int64   `json:"discount_price"`
	CategoryName    string  `json:"category_name" db:"category_name"`
	Ratings         float32 `json:"ratings" db:"ratings"`
	PhotoUrl        string  `json:"photo_url" db:"photo_url"`
	ProductIDString string  `json:"-" db:"product_id_string"`
}

type ManageWishlist struct {
	ProductID    string `json:"product_id"`
	IsWishlisted bool   `json:"is_wishlisted"`
}
