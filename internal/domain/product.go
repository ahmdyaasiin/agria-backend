package domain

type Product struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Quantity    uint   `db:"quantity"`
	Price       int64  `db:"price"`
	CreatedAt   int64  `db:"created_at"`
	UpdatedAt   int64  `db:"updated_at"`
	CategoryID  string `db:"category_id"`
}

func (e Product) TableName() string {
	return "products"
}
