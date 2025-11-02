package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"

	"mini-project-ostore/internal/domain" // Import the domain package
	"mini-project-ostore/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StoreHandler struct {
	storeUC usecase.StoreUseCase
}

func NewStoreHandler(storeUC usecase.StoreUseCase) *StoreHandler {
	return &StoreHandler{storeUC: storeUC}
}

type CreateStoreRequest struct {
	UserID      uint   `form:"user_id" binding:"required"`
	Name        string `form:"name" binding:"required,min=3,max=100"`
	Description string `form:"description"`
	Address     string `form:"address"`
	Phone       string `form:"phone" binding:"omitempty"`
}

type UpdateStoreRequest struct {
	Name        string                `form:"name" binding:"omitempty,min=3,max=100"`
	Description string                `form:"description"`
	Address     string                `form:"address"`
	Phone       string                `form:"phone" binding:"omitempty"`
	PhotoProfile *multipart.FileHeader `form:"photo_profile"` // Field for file upload
}

// PaginatedStoreResponse defines the structure for a paginated list of stores.
type PaginatedStoreResponse struct {
	Stores     []domain.Store `json:"stores"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalCount int64          `json:"total_count"`
	TotalPages int            `json:"total_pages"`
}

// CreateStore handles the creation of a new store.
func (h *StoreHandler) CreateStore(c *gin.Context) {
	var req CreateStoreRequest
	// Changed to c.ShouldBind to handle form-data if files are to be included in CreateStore
	// For now, it's JSON, so c.ShouldBindJSON is more appropriate if no file upload is expected here.
	// Keeping c.ShouldBind for consistency with UpdateStore if future changes involve files.
	if err := c.ShouldBind(&req); err != nil { 
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

// GetStores retrieves all stores with pagination and filtering.
func (h *StoreHandler) GetStores(c *gin.Context) {
	var filter domain.StoreFilter

	// Parse pagination parameters
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			filter.Page = page
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	filter.SetDefaults() // Apply default page and limit if not set

	// Parse filtering parameters
	filter.Search = c.Query("search")

	stores, totalCount, err := h.storeUC.GetStores(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := 0
	if filter.Limit > 0 {
		totalPages = int((totalCount + int64(filter.Limit) - 1) / int64(filter.Limit))
	}

	c.JSON(http.StatusOK, domain.StandardPaginatedResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Data: PaginatedStoreResponse{
			Stores:     stores,
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: totalCount,
			TotalPages: totalPages,
		},
	})
}

// GetUserStores retrieves all stores for a specific user.
// This handler currently expects user ID from path parameter, but in a protected route
// it might come from the authenticated user context.
func (h *StoreHandler) GetUserStores(c *gin.Context) {
	userIDStr := c.Param("id") 
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

// UpdateStore updates an existing store, including handling profile photo upload.
func (h *StoreHandler) UpdateStore(c *gin.Context) {
	storeIDStr := c.Param("id_toko")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	var req UpdateStoreRequest
	// Use c.ShouldBind to handle form-data which includes both text fields and file uploads
	if err := c.ShouldBind(&req); err != nil {
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

	// Handle photo profile upload
	if req.PhotoProfile != nil {
		// Limit upload size to 8MB
		c.Request.ParseMultipartForm(8 << 20) // 8MB

		file := req.PhotoProfile
		if file.Size > (8 << 20) { // 8MB
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file size exceeds 8MB limit"})
			return
		}

		// Generate a unique filename
		extension := filepath.Ext(file.Filename)
		newFileName := uuid.New().String() + extension
		filePath := filepath.Join("uploads", newFileName)

		// Save the file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save photo profile: %v", err)})
			return
		}
		store.PhotoProfile = "/" + filePath // Store the path relative to the server root
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