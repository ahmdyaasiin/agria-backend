package response

type MyCart struct {
	AvailableProducts   []CartProducts `json:"available_products"`
	UnavailableProducts []CartProducts `json:"unavailable_products"`
}

type CartProducts struct {
	ProductID       string `json:"product_id" db:"product_id"`
	ID              string `json:"id" db:"id"`
	Name            string `json:"name" db:"name"`
	Price           int64  `json:"price" db:"price"`
	DiscountPrice   int64  `json:"discount_price"`
	Quantity        int32  `json:"quantity" db:"quantity"`
	PhotoUrl        string `json:"photo_url" db:"photo_url"`
	ProductIDString string `json:"-" db:"product_id_string"`
}

type ManageCart struct {
	ProductID string `json:"product_id"`
	Quantity  uint   `json:"quantity"`
}
