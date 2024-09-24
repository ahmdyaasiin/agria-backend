package request

type ManageCart struct {
	ProductID string `json:"product_id"`
	Quantity  uint   `json:"quantity"`
}
