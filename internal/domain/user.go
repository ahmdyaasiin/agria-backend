package domain

type User struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Username    string `json:"username" db:"username"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Status      string `json:"status" db:"status"`
	IsGoogle    bool   `json:"is_google" db:"is_google"`
	IsFacebook  bool   `json:"is_facebook" db:"is_facebook"`
	PhotoUrl    string `json:"photo_url" db:"photo_url"`
	CreatedAt   int64  `json:"created_at" db:"created_at"`
	UpdatedAt   int64  `json:"updated_at" db:"updated_at"`
}

func (e User) TableName() string {
	return "users"
}
