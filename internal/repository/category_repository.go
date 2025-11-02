package repository

import (
	"mini-project-ostore/internal/domain"

	"gorm.io/gorm"
)

// CategoryRepository defines the interface for category data operations.
type CategoryRepository interface {
	Create(category *domain.Category) error
	FindByID(id uint) (*domain.Category, error)
	Update(category *domain.Category) error
	Delete(id uint) error
	FindAll(filter domain.CategoryFilter) ([]domain.Category, int64, error)
}

// categoryRepository implements the CategoryRepository interface.
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new instance of CategoryRepository.
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// Create a new category in the database.
func (r *categoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

// FindByID retrieves a category by its ID.
func (r *categoryRepository) FindByID(id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	return &category, err
}

// Update an existing category in the database.
func (r *categoryRepository) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

// Delete a category by its ID (soft delete).
func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, id).Error
}

// FindAll retrieves all categories from the database with pagination and filtering.
func (r *categoryRepository) FindAll(filter domain.CategoryFilter) ([]domain.Category, int64, error) {
	var categories []domain.Category
	query := r.db.Model(&domain.Category{})

	// Apply search filter
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
	}

	// Get total count before applying pagination
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.Limit
	query = query.Limit(filter.Limit).Offset(offset)

	err := query.Find(&categories).Error
	return categories, totalCount, err
}