package response

type Homepage struct {
	//
	UserDetails   UserDetails           `json:"user_details"`
	PropertyPromo PropertyPromoHomePage `json:"property_promo"`
	ProductsPromo ProductsPromoHomePage `json:"products_promo"`
	Properties    GetAllProperties      `json:"properties"`
	Products      []GetProduct          `json:"products"`
	Educations    []EducationsHomePage  `json:"educations"`
}

type Market struct {
	//
	UserDetails   UserDetails           `json:"user_details"`
	PropertyPromo PropertyPromoHomePage `json:"property_promo"`
	ProductsPromo ProductsPromoHomePage `json:"products_promo"`
	Properties    PropertyMarket        `json:"properties"`
	Products      ProductMarket         `json:"products"`
}

type PropertyMarket struct {
	//
	Province   string          `json:"province"`
	Provinces  []string        `json:"provinces"`
	Data       []GetProperties `json:"data"`
	Pagination Pagination      `json:"pagination"`
}

type ProductMarket struct {
	//
	Products   []GetProduct `json:"products"`
	Pagination Pagination   `json:"pagination"`
}

type PropertyPromoHomePage struct {
	//
	TimeLifeInSeconds int64           `json:"time_life_in_seconds"`
	Properties        []GetProperties `json:"properties"`
}

type ProductsPromoHomePage struct {
	//
	TimeLifeInSeconds int64        `json:"time_life_in_seconds"`
	Products          []GetProduct `json:"products"`
}

type EducationsHomePage struct {
	ID             string `json:"id" db:"id"`
	Title          string `json:"title" db:"title"`
	PhotoUrl       string `json:"photo_url" db:"photo_url"`
	PhotoUrlAuthor string `json:"photo_url_author" db:"photo_url_author"`
	NameOfAuthor   string `json:"name_of_author" db:"name"`
	CreatedAt      int64  `json:"created_at" db:"created_at"`
	CountLikes     int64  `json:"count_likes"`
	CountViews     int64  `json:"count_views"`
	CountComments  int64  `json:"count_comments"`
	InWishlist     bool   `json:"in_wishlist"`
}
