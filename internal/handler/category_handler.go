package handler

import (
	"net/http"
	"strconv"

	"mini-project-ostore/internal/domain" // Import the domain package
	"mini-project-ostore/internal/usecase"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryUC usecase.CategoryUseCase
}

func NewCategoryHandler(categoryUC usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{categoryUC: categoryUC}
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,min=3,max=100"`
	Description string `json:"description"`
}

// PaginatedCategoryResponse defines the structure for a paginated list of categories.
type PaginatedCategoryResponse struct {
	Categories []domain.Category `json:"categories"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalCount int64             `json:"total_count"`
	TotalPages int               `json:"total_pages"`
}

// CreateCategory handles the creation of a new category.
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.categoryUC.Create(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully", "category_id": category.ID})
}

// GetCategories retrieves all categories with pagination and filtering.
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var filter domain.CategoryFilter

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

	categories, totalCount, err := h.categoryUC.GetCategories(filter)
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
		Data: PaginatedCategoryResponse{
			Categories: categories,
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: totalCount,
			TotalPages: totalPages,
		},
	})
}

// GetCategoryByID retrieves a single category by its ID.
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := h.categoryUC.GetByID(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// UpdateCategory updates an existing category.
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := &domain.Category{ID: uint(categoryID)}
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := h.categoryUC.Update(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

// DeleteCategory deletes a category by its ID.
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := h.categoryUC.Delete(uint(categoryID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}