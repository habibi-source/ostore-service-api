package domain

// TransactionFilter represents the filters and pagination parameters for transaction queries.
type TransactionFilter struct {
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
	UserID      uint   `json:"user_id"`      // Filter by user who made the transaction
	StoreID     uint   `json:"store_id"`     // Filter by store involved in the transaction
	Status      string `json:"status"`       // Filter by transaction status
	PaymentMethod string `json:"payment_method"` // Filter by payment method
}

// SetDefaults sets default values for pagination if not provided.
func (f *TransactionFilter) SetDefaults() {
	if f.Page <= 0 {
		f.Page = DefaultPage
	}
	if f.Limit <= 0 {
		f.Limit = DefaultLimit
	}
}
