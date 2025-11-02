package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"mini-project-ostore/internal/domain"
	"net/http"
	"time"
)

// CityAPIRepository defines the interface for fetching city data from an external API.
type CityAPIRepository interface {
	GetCitiesByProvinceID(provinceID string) ([]domain.City, error)
}

// cityAPIRepository implements the CityAPIRepository interface by
// fetching data from the EMSIFA API.
type cityAPIRepository struct {
	client *http.Client
}

// NewCityAPIRepository creates a new instance of CityAPIRepository.
func NewCityAPIRepository() CityAPIRepository {
	return &cityAPIRepository{
		client: &http.Client{
			Timeout: 10 * time.Second, // Set a timeout for HTTP requests
		},
	}
}

// GetCitiesByProvinceID fetches cities for a given province from the EMSIFA API.
func (r *cityAPIRepository) GetCitiesByProvinceID(provinceID string) ([]domain.City, error) {
	resp, err := r.client.Get(fmt.Sprintf("%s/regencies/%s.json", emsifaBaseURL, provinceID))
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

	var cities []domain.City
	if err := json.Unmarshal(body, &cities); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cities data: %w", err)
	}

	return cities, nil
}