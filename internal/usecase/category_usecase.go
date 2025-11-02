package usecase

import (
	"errors"
	"mini-project-ostore/internal/domain"
	"mini-project-ostore/internal/repository"
)

// CategoryUseCase defines the interface for category-related business logic.
type CategoryUseCase interface {
	Create(category *domain.Category) error
	GetByID(id uint) (*domain.Category, error)
	Update(category *domain.Category) error
	Delete(id uint) error
	GetCategories(filter domain.CategoryFilter) ([]domain.Category, int64, error)
}

// categoryUseCase implements the CategoryUseCase interface.
type categoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryUseCase creates a new instance of CategoryUseCase.
func NewCategoryUseCase(categoryRepo repository.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{categoryRepo: categoryRepo}
}

// Create a new category.
func (uc *categoryUseCase) Create(category *domain.Category) error {
	// Optional: Add business logic validation here.
	// For example, check for duplicate category names before creating.
	return uc.categoryRepo.Create(category)
}

// GetByID retrieves a category by its ID.
func (uc *categoryUseCase) GetByID(id uint) (*domain.Category, error) {
	return uc.categoryRepo.FindByID(id)
}

// Update an existing category.
func (uc *categoryUseCase) Update(category *domain.Category) error {
	// First, check if the category exists.
	existingCategory, err := uc.categoryRepo.FindByID(category.ID)
	if err != nil {
		return errors.New("category not found")
	}

	// Update only the mutable fields.
	existingCategory.Name = category.Name
	existingCategory.Description = category.Description

	return uc.categoryRepo.Update(existingCategory)
}

// Delete a category by its ID.
func (uc *categoryUseCase) Delete(id uint) error {
	// First, check if the category exists.
	_, err := uc.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}
	return uc.categoryRepo.Delete(id)
}

// GetCategories retrieves all categories with pagination and filtering.
func (uc *categoryUseCase) GetCategories(filter domain.CategoryFilter) ([]domain.Category, int64, error) {
	return uc.categoryRepo.FindAll(filter)
}