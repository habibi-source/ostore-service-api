package domain

// StoreFilter represents the filters and pagination parameters for store queries.
type StoreFilter struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Search string `json:"search"` // Adding search for future potential filtering
}

// SetDefaults sets default values for pagination if not provided.
func (f *StoreFilter) SetDefaults() {
	if f.Page <= 0 {
		f.Page = DefaultPage
	}
	if f.Limit <= 0 {
		f.Limit = DefaultLimit
	}
}