package domain

type Address struct {
	//
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	Address     string  `db:"address"`
	City        string  `db:"city"`
	State       string  `db:"state"`
	PostalCode  string  `db:"postal_code"`
	Latitude    float64 `db:"latitude"`
	Longitude   float64 `db:"longitude"`
	IsPrimary   bool    `db:"is_primary"`
	PhoneNumber string  `db:"phone_number"`
	CreatedAt   int64   `db:"created_at"`
	UpdatedAt   int64   `db:"updated_at"`
	UserID      string  `db:"user_id"`
}

func (e Address) TableName() string {
	return "addresses"
}
