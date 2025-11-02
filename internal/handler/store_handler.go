package handler

import (
	"net/http"
	"strconv"

	"mini-project-ostore/internal/domain"
	"mini-project-ostore/internal/usecase"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeUC usecase.StoreUseCase
}

func NewStoreHandler(storeUC usecase.StoreUseCase) *StoreHandler {
	return &StoreHandler{storeUC: storeUC}
}

type CreateStoreRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Phone       string `json:"phone" binding:"omitempty"`
}

type UpdateStoreRequest struct {
	Name        string `json:"name" binding:"omitempty,min=3,max=100"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Phone       string `json:"phone" binding:"omitempty"`
}

// CreateStore handles the creation of a new store.
func (h *StoreHandler) CreateStore(c *gin.Context) {
	var req CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store := &domain.Store{
		UserID:      req.UserID,
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Phone:       req.Phone,
	}

	if err := h.storeUC.Create(store); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Store created successfully", "store_id": store.ID})
}

// GetStores retrieves all stores.
func (h *StoreHandler) GetStores(c *gin.Context) {
	stores, err := h.storeUC.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stores)
}

// GetUserStores retrieves all stores for a specific user.
func (h *StoreHandler) GetUserStores(c *gin.Context) {
	userIDStr := c.Param("id") // Assuming the user ID is passed as a path parameter, e.g., /users/:id/stores
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stores, err := h.storeUC.GetByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stores)
}

// GetStoreByID retrieves a single store by its ID.
func (h *StoreHandler) GetStoreByID(c *gin.Context) {
	storeIDStr := c.Param("id_toko")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	store, err := h.storeUC.GetByID(uint(storeID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, store)
}

// UpdateStore updates an existing store.
func (h *StoreHandler) UpdateStore(c *gin.Context) {
	storeIDStr := c.Param("id_toko")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	var req UpdateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store := &domain.Store{ID: uint(storeID)}
	if req.Name != "" {
		store.Name = req.Name
	}
	if req.Description != "" {
		store.Description = req.Description
	}
	if req.Address != "" {
		store.Address = req.Address
	}
	if req.Phone != "" {
		store.Phone = req.Phone
	}

	if err := h.storeUC.Update(store); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Store updated successfully"})
}

// DeleteStore deletes a store by its ID.
func (h *StoreHandler) DeleteStore(c *gin.Context) {
	storeIDStr := c.Param("id_toko")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	if err := h.storeUC.Delete(uint(storeID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Store deleted successfully"})
}
// GetMyStore retrieves the store of the currently authenticated user.
func (h *StoreHandler) GetMyStore(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	store, err := h.storeUC.GetByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Succeed to GET data",
		"data":    store,
	})
}
