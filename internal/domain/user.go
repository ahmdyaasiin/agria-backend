package domain

type User struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Username    string `db:"username"`
	Email       string `db:"email"`
	Password    string `db:"password"`
	PhoneNumber string `db:"phone_number"`
	Status      string `db:"status"`
	IsGoogle    bool   `db:"is_google"`
	IsFacebook  bool   `db:"is_facebook"`
	UrlPhoto    string `db:"url_photo"`
	CreatedAt   int64  `db:"created_at"`
	UpdatedAt   int64  `db:"updated_at"`
}

func (e User) TableName() string {
	return "users"
}
