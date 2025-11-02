package handler

import (
	"net/http"

	"mini-project-ostore/internal/usecase"

	"github.com/gin-gonic/gin"
)

type RegionHandler struct {
	regionUC usecase.RegionUseCase
}

func NewRegionHandler(regionUC usecase.RegionUseCase) *RegionHandler {
	return &RegionHandler{regionUC: regionUC}
}

// GetProvinces retrieves all provinces.
func (h *RegionHandler) GetProvinces(c *gin.Context) {
	provinces, err := h.regionUC.GetAllProvinces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, provinces)
}

// GetCities retrieves cities for a given province ID.
func (h *RegionHandler) GetCities(c *gin.Context) {
	provinceID := c.Param("provinceID")
	if provinceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Province ID is required"})
		return
	}

	cities, err := h.regionUC.GetCitiesByProvinceID(provinceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cities)
}

// GetSubdistricts retrieves subdistricts for a given city ID.
func (h *RegionHandler) GetSubdistricts(c *gin.Context) {
	cityID := c.Param("cityID")
	if cityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City ID is required"})
		return
	}

	subdistricts, err := h.regionUC.GetSubdistrictsByCityID(cityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subdistricts)
}
