package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"mini-project-ostore/internal/domain"
	"net/http"
	"time"
)

const emsifaBaseURL = "https://www.emsifa.com/api-wilayah-indonesia/api"

// ProvinceAPIRepository defines the interface for fetching province data from an external API.
type ProvinceAPIRepository interface {
	GetAllProvinces() ([]domain.Province, error)
}

// provinceAPIRepository implements the ProvinceAPIRepository interface by
// fetching data from the EMSIFA API.
type provinceAPIRepository struct {
	client *http.Client
}

// NewProvinceAPIRepository creates a new instance of ProvinceAPIRepository.
func NewProvinceAPIRepository() ProvinceAPIRepository {
	return &provinceAPIRepository{
		client: &http.Client{
			Timeout: 10 * time.Second, // Set a timeout for HTTP requests
		},
	}
}

// GetAllProvinces fetches all provinces from the EMSIFA API.
func (r *provinceAPIRepository) GetAllProvinces() ([]domain.Province, error) {
	resp, err := r.client.Get(fmt.Sprintf("%s/provinces.json", emsifaBaseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to make request to EMSIFA API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code from EMSIFA API: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var provinces []domain.Province
	if err := json.Unmarshal(body, &provinces); err != nil {
		return nil, fmt.Errorf("failed to unmarshal provinces data: %w", err)
	}

	return provinces, nil
}