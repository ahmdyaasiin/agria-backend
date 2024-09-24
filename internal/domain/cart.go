package domain

type Cart struct {
	//
	ID        string `db:"id"`
	Quantity  uint   `db:"quantity"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
	UserID    string `db:"user_id"`
	ProductID string `db:"product_id"`
}

func (e Cart) TableName() string {
	return "carts"
}
