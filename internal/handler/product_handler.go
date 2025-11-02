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

// CreateProduct handles the creation of a new product.
func (h *ProductHandler) CreateProduct(c *gin.Context) {
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

	if err := h.productUC.Create(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product_id": product.ID})
}

// GetProducts retrieves all products.
func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.productUC.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetUserProducts retrieves all products for stores owned by a specific user.
func (h *ProductHandler) GetUserProducts(c *gin.Context) {
	userIDStr := c.Param("id") // Assuming the user ID is passed as a path parameter, e.g., /users/:id/products
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	products, err := h.productUC.GetUserProducts(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
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
	// Only update fields if they are explicitly provided in the request
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
	if req.Stock != 0 { // This might be an issue if 0 is a valid stock update
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

	if err := h.productUC.Update(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// DeleteProduct deletes a product by its ID.
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.productUC.Delete(uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}