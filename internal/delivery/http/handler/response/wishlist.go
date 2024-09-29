package response

type UserDetails struct {
	//
	IsLoggedIn         bool   `json:"is_logged_in"`
	CountCarts         int    `json:"count_carts"`
	CountNotifications int    `json:"count_notifications"`
	PhotoProfile       string `json:"photo_profile"`
}

type ProductWishlist struct {
	UserDetails UserDetails  `json:"user_details"`
	Products    []MyWishlist `json:"products"`
	Pagination  Pagination   `json:"pagination"`
}

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

type PropertiesWishlist struct {
	UserDetails UserDetails            `json:"user_details"`
	Properties  []MyWishlistProperties `json:"properties"`
	Pagination  Pagination             `json:"pagination"`
}

type MyWishlistProperties struct {
	//
	ID                string  `json:"id" db:"id"`
	Name              string  `json:"name" db:"name"`
	Price             int64   `json:"price" db:"price"`
	DiscountPrice     int64   `json:"discount_price"`
	CategoryName      string  `json:"category_name" db:"category_name"`
	Ratings           float32 `json:"star" db:"ratings"`
	PhotoUrl          string  `json:"photo_url" db:"photo_url"`
	CertificationType string  `json:"ownership_type" db:"certification_type"`
	Width             int     `json:"width" db:"width"`
	City              string  `json:"city" db:"city"`
	State             string  `json:"state" db:"state"`
	IsWishlisted      bool    `json:"in_wishlist" db:"in_wishlist"`
	ProductIDString   string  `json:"-" db:"product_id_string"`
}

type ManageWishlist struct {
	ProductID    string `json:"product_id"`
	IsWishlisted bool   `json:"in_wishlist"`
}

type ManageWishlistProperties struct {
	PropertiesID string `json:"property_id"`
	IsWishlisted bool   `json:"is_wishlisted"`
}
