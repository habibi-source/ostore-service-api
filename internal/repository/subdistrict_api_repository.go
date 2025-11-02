package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"mini-project-ostore/internal/domain"
	"net/http"
	"time"
)

// SubdistrictAPIRepository defines the interface for fetching subdistrict data from an external API.
type SubdistrictAPIRepository interface {
	GetSubdistrictsByCityID(cityID string) ([]domain.Subdistrict, error)
}

// subdistrictAPIRepository implements the SubdistrictAPIRepository interface by
// fetching data from the EMSIFA API.
type subdistrictAPIRepository struct {
	client *http.Client
}

// NewSubdistrictAPIRepository creates a new instance of SubdistrictAPIRepository.
func NewSubdistrictAPIRepository() SubdistrictAPIRepository {
	return &subdistrictAPIRepository{
		client: &http.Client{
			Timeout: 10 * time.Second, // Set a timeout for HTTP requests
		},
	}
}

// GetSubdistrictsByCityID fetches subdistricts for a given city (regency) from the EMSIFA API.
func (r *subdistrictAPIRepository) GetSubdistrictsByCityID(cityID string) ([]domain.Subdistrict, error) {
	resp, err := r.client.Get(fmt.Sprintf("%s/districts/%s.json", emsifaBaseURL, cityID))
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

	var subdistricts []domain.Subdistrict
	if err := json.Unmarshal(body, &subdistricts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal subdistricts data: %w", err)
	}

	return subdistricts, nil
}