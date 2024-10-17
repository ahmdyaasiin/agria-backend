package response

type EducationDetails struct {
	UserDetails    UserDetails   `json:"user_details"`
	Data           EducationData `json:"data"`
	RelatedArticle struct {
		Data       []EducationCard `json:"data"`
		Pagination Pagination      `json:"pagination"`
	} `json:"related_article"`
}

type EducationData struct {
	ID              string `json:"id" db:"id"`
	Title           string `json:"title" db:"title"`
	PhotoUrl        string `json:"photo_url" db:"photo_url"`
	Author          string `json:"author" db:"author"`
	PhotoUrlAuthor  string `json:"photo_url_author" db:"photo_url_author"`
	Content         string `json:"content" db:"content"`
	CreatedAt       int64  `json:"-" db:"created_at"`
	CreatedAtString string `json:"created_at"`
}
