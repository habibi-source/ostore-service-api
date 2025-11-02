package usecase

import (
	"mini-project-ostore/internal/domain"
	"mini-project-ostore/internal/repository"
)

// RegionUseCase defines the interface for region-related business logic.
type RegionUseCase interface {
	GetAllProvinces() ([]domain.Province, error)
	GetCitiesByProvinceID(provinceID string) ([]domain.City, error)
	GetSubdistrictsByCityID(cityID string) ([]domain.Subdistrict, error)
}

// regionUseCase implements the RegionUseCase interface.
type regionUseCase struct {
	provinceAPIRepo    repository.ProvinceAPIRepository
	cityAPIRepo        repository.CityAPIRepository
	subdistrictAPIRepo repository.SubdistrictAPIRepository
}

// NewRegionUseCase creates a new instance of RegionUseCase.
func NewRegionUseCase(
	provinceAPIRepo repository.ProvinceAPIRepository,
	cityAPIRepo repository.CityAPIRepository,
	subdistrictAPIRepo repository.SubdistrictAPIRepository,
) RegionUseCase {
	return &regionUseCase{
		provinceAPIRepo:    provinceAPIRepo,
		cityAPIRepo:        cityAPIRepo,
		subdistrictAPIRepo: subdistrictAPIRepo,
	}
}

// GetAllProvinces retrieves all provinces using the API repository.
func (uc *regionUseCase) GetAllProvinces() ([]domain.Province, error) {
	return uc.provinceAPIRepo.GetAllProvinces()
}

// GetCitiesByProvinceID retrieves cities for a given province ID using the API repository.
func (uc *regionUseCase) GetCitiesByProvinceID(provinceID string) ([]domain.City, error) {
	return uc.cityAPIRepo.GetCitiesByProvinceID(provinceID)
}

// GetSubdistrictsByCityID retrieves subdistricts for a given city ID using the API repository.
func (uc *regionUseCase) GetSubdistrictsByCityID(cityID string) ([]domain.Subdistrict, error) {
	return uc.subdistrictAPIRepo.GetSubdistrictsByCityID(cityID)
}
