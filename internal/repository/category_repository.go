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
	FindAll() ([]domain.Category, error)
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

// FindAll retrieves all categories from the database.
func (r *categoryRepository) FindAll() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.Find(&categories).Error
	return categories, err
}