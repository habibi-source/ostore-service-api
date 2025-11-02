package domain

// ProductFilter represents the filters and pagination parameters for product queries.
type ProductFilter struct {
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	Search     string  `json:"search"`
	CategoryID uint    `json:"category_id"`
	MinPrice   float64 `json:"min_price"`
	MaxPrice   float64 `json:"max_price"`
	StoreID    uint    `json:"store_id"`
	UserID     uint    `json:"user_id"`
}

// Default values for pagination
const (
	DefaultPage  = 1
	DefaultLimit = 10
)

// SetDefaults sets default values for pagination if not provided.
func (f *ProductFilter) SetDefaults() {
	if f.Page <= 0 {
		f.Page = DefaultPage
	}
	if f.Limit <= 0 {
		f.Limit = DefaultLimit
	}
}