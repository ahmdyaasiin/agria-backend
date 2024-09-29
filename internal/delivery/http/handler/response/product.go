package response

type GetProductWithPagination struct {
	Products []GetProduct `json:"products"`
	Pagination
}

type GetPropertiesWithPagination struct {
	Properties []GetProperties `json:"properties"`
	Pagination
}

type GetProperties struct {
	//
	ID                string  `json:"id" db:"id"`
	Name              string  `json:"name" db:"name"`
	CategoryName      string  `json:"category_name" db:"category_name"`
	City              string  `json:"city" db:"city"`
	Price             int64   `json:"price" db:"price"`
	Width             int     `json:"width" db:"width"`
	CertificationType string  `json:"certification_type" db:"certification_type"`
	PhotoUrl          string  `json:"photo_url" db:"photo_url"`
	Ratings           float32 `json:"ratings" db:"ratings"`
	IsWishlisted      bool    `json:"is_wishlisted" db:"is_wishlisted"`
}

type GetProduct struct {
	ID              string  `json:"id" db:"id"`
	Name            string  `json:"name" db:"name"`
	Price           int64   `json:"price" db:"price"`
	DiscountPrice   int64   `json:"discount_price"`
	CategoryName    string  `json:"category_name" db:"category_name"`
	Ratings         float32 `json:"ratings" db:"ratings"`
	IsWishlisted    bool    `json:"is_wishlisted" db:"is_wishlisted"`
	PhotoUrl        string  `json:"photo_url" db:"photo_url"`
	ProductIDString string  `json:"-" db:"product_id_string"`

	// unused key
	Description          string `json:"-" db:"description"`
	Quantity             uint   `json:"-" db:"quantity"`
	UnitWeight           int32  `json:"-" db:"unit_weight"`
	ShelfLife            int32  `json:"-" db:"shelf_life"`
	OrganicCertification string `json:"-" db:"organic_certification"`
	CreatedAt            int64  `json:"-" db:"created_at"`
	UpdatedAt            int64  `json:"-" db:"updated_at"`
	CategoryID           string `json:"-" db:"category_id"`
}

type GetProductDetails struct {
	//
	ID                   string   `json:"id" db:"id"`
	Name                 string   `json:"name" db:"name"`
	Description          string   `json:"description"`
	Quantity             uint     `json:"quantity" db:"quantity"`
	Price                int64    `json:"price" db:"price"`
	UnitWeight           int32    `json:"unit_weight" db:"unit_weight"`
	ShelfLife            int32    `json:"shelf_life" db:"shelf_life"`
	OrganicCertification string   `json:"organic_certification" db:"organic_certification"`
	PhotoUrls            []string `json:"photo_urls"`
	TimeRange            string   `json:"time_range"`
	PriceRange           string   `json:"price_range"`
	CategoryName         string   `json:"category_name" db:"category_name"`
	Ratings              float32  `json:"ratings" db:"ratings"`
	IsWishlisted         bool     `json:"is_wishlisted" db:"is_wishlisted"`
	ReviewsCount         uint     `json:"reviews_count" db:"reviews_count"`
	Reviews              []Review `json:"reviews"`

	// unused key
	CreatedAt  int64  `json:"-" db:"created_at"`
	UpdatedAt  int64  `json:"-" db:"updated_at"`
	CategoryID string `json:"-" db:"category_id"`
}

type GetPropertyDetails struct {
	//
	ID            string `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	Description   string `json:"description" db:"description"`
	Price         int64  `json:"price" db:"price"`
	CategoryName  string `json:"category_name" db:"category_name"`
	Width         int    `json:"width" db:"width"`
	Province      string `json:"province" db:"province"`
	City          string `json:"city" db:"city"`
	OwnershipType string `json:"ownership_type" db:"ownership_type"`
	NameOfOwner   string `json:"name_of_owner" db:"name_of_owner"`
	PhotoUrl      string `json:"photo_url" db:"photo_url"`

	//
	DiscountPrice     int64                `json:"discount_price" db:"discount_price"`
	InWishlist        bool                 `json:"in_wishlist" db:"in_wishlist"`
	PhotoUrls         []string             `json:"photo_urls" db:"photo_urls"`
	RatingsAndReviews []RatingProperty     `json:"ratings_and_reviews"`
	Discuss           []PropertyDiscuss    `json:"discuss"`
	Highlights        []PropertyHighlights `json:"highlights"`
}

type RatingAndReviewsProperty struct {
	//
	CountRatings       int              `json:"count_ratings"`
	CountStarBreakDown []int            `json:"count_star_break_down"`
	Data               []RatingProperty `json:"data"`
}

type RatingProperty struct {
	ID              string   `json:"id" db:"id"`
	Name            string   `json:"name" db:"name"`
	PhotoUrl        string   `json:"photo_url" db:"photo_url"`
	Content         string   `json:"content" db:"content"`
	CountHelpful    int      `json:"count_helpful" db:"count_helpful"`
	IsHelpful       bool     `json:"is_helpful" db:"is_helpful"`
	PhotoUrls       []string `json:"photo_urls" db:"photo_urls"`
	PhotoUrlsString string   `json:"-" db:"photo_urls_string"`
}

type PropertyDiscuss struct {
	//
	ID            string                   `json:"id" db:"id"`
	Name          string                   `json:"name" db:"name"`
	Content       string                   `json:"content" db:"content"`
	PhotoUrl      string                   `json:"photo_url" db:"photo_url"`
	Answers       []PropertyDiscussReplies `json:"answers"`
	AnswersString string                   `json:"-" db:"answers_string"`
}

type PropertyDiscussReplies struct {
	//
	ID       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Content  string `json:"content" db:"content"`
	PhotoUrl string `json:"photo_url" db:"photo_url"`
	IsOwner  bool   `json:"is_owner"`
}

type PropertyHighlights struct {
	//
	Name     string `json:"name" db:"name"`
	PhotoUrl string `json:"photo_url" db:"photo_url"`
}

type ReviewDetails struct {
	Reviews         []Review `json:"reviews"`
	RatingBreakdown []int64  `json:"rating_breakdown"`
	Pagination
}

type RatingBreakdown struct {
	Star  int   `json:"star" db:"star"`
	Total int64 `json:"total" db:"total"`
}

type Review struct {
	ID                    string   `json:"id" db:"id"`
	Name                  string   `json:"name" db:"name"`
	PhotoUrl              string   `json:"photo_url" db:"photo_url"`
	Star                  int      `json:"star" db:"star"`
	Content               string   `json:"content" db:"content"`
	PhotoReviewUrlsString string   `json:"-" db:"photo_reviews_urls_string"`
	PhotoReviewUrls       []string `json:"photo_review_urls"`
	HelpfulCount          int      `json:"helpful_count" db:"helpful_count"`
	IsReviewHelpful       bool     `json:"is_review_helpful" db:"is_review_helpful"`
	CreatedAt             int64    `json:"created_at" db:"created_at"`

	// unused key
	UpdatedAt         int64  `json:"-" db:"updated_at"`
	UserID            string `json:"-" db:"user_id"`
	TransactionItemID string `json:"-" db:"transaction_item_id"`
}
