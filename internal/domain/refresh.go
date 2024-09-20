package domain

type Refresh struct {
	//
	ID            string `db:"id"`
	Token         string `db:"token"`
	CreatedAt     int64  `db:"created_at"`
	LastRefreshAt int64  `db:"last_refresh_at"`
	UserID        string `db:"user_id"`
}

func (e Refresh) TableName() string {
	return "refreshes"
}
