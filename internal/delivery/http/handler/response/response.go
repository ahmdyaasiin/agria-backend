package response

type Pagination struct {
	Page       int   `json:"page"`
	TotalItems int64 `json:"count_items"`
	TotalPages int64 `json:"count_pages"`
}
