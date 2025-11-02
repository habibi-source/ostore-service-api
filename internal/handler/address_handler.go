package handler

import (
	"net/http"
	"strconv"

	"mini-project-ostore/internal/domain"
	"mini-project-ostore/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	addressUC usecase.AddressUseCase
}

func NewAddressHandler(addressUC usecase.AddressUseCase) *AddressHandler {
	return &AddressHandler{addressUC: addressUC}
}

type CreateAddressRequest struct {
	UserID        uint   `json:"user_id" binding:"required"`
	Label         string `json:"label" binding:"required,min=3,max=50"`
	ReceiverName  string `json:"receiver_name" binding:"required,min=3,max=100"`
	Phone         string `json:"phone" binding:"required"`
	ProvinceID    uint `json:"province_id" binding:"required"`
	CityID        uint`json:"city_id" binding:"required"`
	SubDistrictID uint `json:"subdistrict_id" binding:"required"`
	Detail        string `json:"detail" binding:"required"`
	PostalCode    string `json:"postal_code" binding:"required"`
	IsPrimary     *bool  `json:"is_primary"` // Use pointer to differentiate between false and not provided
}

type UpdateAddressRequest struct {
	Label         string `json:"label" binding:"omitempty,min=3,max=50"`
	ReceiverName  string `json:"receiver_name" binding:"omitempty,min=3,max=100"`
	Phone         string `json:"phone" binding:"omitempty"`
	ProvinceID    uint `json:"province_id" binding:"omitempty"`
	CityID        uint `json:"city_id" binding:"omitempty"`
	SubDistrictID uint `json:"subdistrict_id" binding:"omitempty"`
	Detail        string `json:"detail" binding:"omitempty"`
	PostalCode    string `json:"postal_code" binding:"omitempty"`
	IsPrimary     *bool  `json:"is_primary"` // Use pointer to differentiate between false and not provided
}

// CreateAddress handles the creation of a new address.
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address := &domain.Address{
		UserID:        req.UserID,
		Label:         req.Label,
		ReceiverName:  req.ReceiverName,
		Phone:         req.Phone,
		ProvinceID:    req.ProvinceID,
		CityID:        req.CityID,
		SubDistrictID: req.SubDistrictID,
		Detail:        req.Detail,
		PostalCode:    req.PostalCode,
		IsPrimary:     (func() bool {
			if req.IsPrimary != nil {
				return *req.IsPrimary
			}
			return false // Default to false if not provided, consistent with gorm default
		})(),
	}

	if err := h.addressUC.Create(address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Address created successfully", "address_id": address.ID})
}

// GetUserAddresses retrieves all addresses for a specific user.
func (h *AddressHandler) GetUserAddresses(c *gin.Context) {
	userIDStr := c.Param("id") // Assuming user ID is passed as a path parameter, e.g., /users/:id/addresses
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	addresses, err := h.addressUC.GetByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, addresses)
}

// GetAddress retrieves a single address by its ID.
func (h *AddressHandler) GetAddress(c *gin.Context) {
	addressIDStr := c.Param("id")
	addressID, err := strconv.ParseUint(addressIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	address, err := h.addressUC.GetByID(uint(addressID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, address)
}

// UpdateAddress updates an existing address.
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	addressIDStr := c.Param("id")
	addressID, err := strconv.ParseUint(addressIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	var req UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address := &domain.Address{ID: uint(addressID)}
	if req.Label != "" {
		address.Label = req.Label
	}
	if req.ReceiverName != "" {
		address.ReceiverName = req.ReceiverName
	}
	if req.Phone != "" {
		address.Phone = req.Phone
	}
	if req.ProvinceID != 0 {
		address.ProvinceID = req.ProvinceID
	}
	if req.CityID != 0 {
		address.CityID = req.CityID
	}
	if req.SubDistrictID != 0 {
		address.SubDistrictID = req.SubDistrictID
	}
	if req.Detail != "" {
		address.Detail = req.Detail
	}
	if req.PostalCode != "" {
		address.PostalCode = req.PostalCode
	}
	if req.IsPrimary != nil {
		address.IsPrimary = *req.IsPrimary
	}

	if err := h.addressUC.Update(address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address updated successfully"})
}

// DeleteAddress deletes an address by its ID.
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	addressIDStr := c.Param("id")
	addressID, err := strconv.ParseUint(addressIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	if err := h.addressUC.Delete(uint(addressID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}
