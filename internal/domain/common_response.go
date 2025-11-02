package domain

// StandardPaginatedResponse defines the standard structure for paginated API responses.
type StandardPaginatedResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"` // This will hold specific paginated data (e.g., PaginatedProductResponse)
}