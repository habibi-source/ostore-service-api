package handler

import (
	"net/http"
	"strconv"

	"mini-project-ostore/internal/domain"
	"mini-project-ostore/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUC usecase.ProductUseCase
}

func NewProductHandler(productUC usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{productUC: productUC}
}

type CreateProductRequest struct {
	StoreID     uint    `json:"store_id" binding:"required"`
	CategoryID  uint    `json:"category_id" binding:"required"`
	SKU         string  `json:"sku" binding:"required,min=3,max=50"`
	Slug        string  `json:"slug" binding:"required,min=3,max=255"`
	Name        string  `json:"name" binding:"required,min=3,max=200"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	Weight      float64 `json:"weight" binding:"omitempty,gt=0"`
	Images      string  `json:"images"` // JSON array of image URLs
	IsAvailable bool    `json:"is_available"`
}

type UpdateProductRequest struct {
	StoreID     uint    `json:"store_id" binding:"omitempty"`
	CategoryID  uint    `json:"category_id" binding:"omitempty"`
	SKU         string  `json:"sku" binding:"omitempty,min=3,max=50"`
	Slug        string  `json:"slug" binding:"omitempty,min=3,max=255"`
	Name        string  `json:"name" binding:"omitempty,min=3,max=200"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       int     `json:"stock" binding:"omitempty,gte=0"`
	Weight      float64 `json:"weight" binding:"omitempty,gt=0"`
	Images      string  `json:"images"` // JSON array of image URLs
	IsAvailable *bool   `json:"is_available"`
}

// PaginatedProductResponse defines the structure for a paginated list of products.
type PaginatedProductResponse struct {
	Products   []domain.Product `json:"products"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalCount int64            `json:"total_count"`
	TotalPages int              `json:"total_pages"`
}

// CreateProduct handles the creation of a new product.
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authenticatedUserID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &domain.Product{
		StoreID:     req.StoreID,
		CategoryID:  req.CategoryID,
		SKU:         req.SKU,
		Slug:        req.Slug,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Weight:      req.Weight,
		Images:      req.Images,
		IsAvailable: req.IsAvailable,
	}

	// Pass userID to Create
	if err := h.productUC.Create(product, authenticatedUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product_id": product.ID})
}

// GetProducts retrieves all products with pagination and filtering.
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var filter domain.ProductFilter

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

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			filter.CategoryID = uint(categoryID)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id format"})
			return
		}
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filter.MinPrice = minPrice
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min_price format"})
			return
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filter.MaxPrice = maxPrice
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max_price format"})
			return
		}
	}

	if storeIDStr := c.Query("store_id"); storeIDStr != "" {
		if storeID, err := strconv.ParseUint(storeIDStr, 10, 32); err == nil {
			filter.StoreID = uint(storeID)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store_id format"})
			return
		}
	}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			filter.UserID = uint(userID)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
			return
		}
	}

	products, totalCount, err := h.productUC.GetProducts(filter)
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
		Data: PaginatedProductResponse{
			Products:   products,
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: totalCount,
			TotalPages: totalPages,
		},
	})
}

// GetUserProducts retrieves all products for stores owned by the authenticated user with pagination and filtering.
func (h *ProductHandler) GetUserProducts(c *gin.Context) {
	// Get authenticated user ID from context
	// This assumes the auth middleware adds "userID" to the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authenticatedUserID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	var filter domain.ProductFilter
	// Set the UserID from the authenticated user
	filter.UserID = authenticatedUserID

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

	// Parse filtering parameters (similar to GetProducts)
	filter.Search = c.Query("search")

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			filter.CategoryID = uint(categoryID)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id format"})
			return
		}
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filter.MinPrice = minPrice
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min_price format"})
			return
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filter.MaxPrice = maxPrice
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max_price format"})
			return
		}
	}

	// StoreID is for filtering products by a specific store.
	// This is different from the user's *own* store, but allows filtering within their products.
	if storeIDStr := c.Query("store_id"); storeIDStr != "" {
		if storeID, err := strconv.ParseUint(storeIDStr, 10, 32); err == nil {
			filter.StoreID = uint(storeID)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store_id format"})
			return
		}
	}

	products, totalCount, err := h.productUC.GetUserProducts(filter)
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
		Data: PaginatedProductResponse{
			Products:   products,
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: totalCount,
			TotalPages: totalPages,
		},
	})
}

// GetProductByID retrieves a single product by its ID.
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.productUC.GetByID(uint(productID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// UpdateProduct updates an existing product.
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authenticatedUserID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &domain.Product{ID: uint(productID)}
	if req.StoreID != 0 {
		product.StoreID = req.StoreID
	}
	if req.CategoryID != 0 {
		product.CategoryID = req.CategoryID
	}
	if req.SKU != "" {
		product.SKU = req.SKU
	}
	if req.Slug != "" {
		product.Slug = req.Slug
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != 0 {
		product.Price = req.Price
	}
	if req.Stock != 0 {
		product.Stock = req.Stock
	}
	if req.Weight != 0 {
		product.Weight = req.Weight
	}
	if req.Images != "" {
		product.Images = req.Images
	}
	if req.IsAvailable != nil {
		product.IsAvailable = *req.IsAvailable
	}

	// Pass userID to Update
	if err := h.productUC.Update(product, authenticatedUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// DeleteProduct deletes a product by its ID.
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authenticatedUserID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Pass both productID and userID to Delete
	if err := h.productUC.Delete(uint(productID), authenticatedUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
