package repository

import (
	"mini-project-ostore/internal/domain"

	"gorm.io/gorm"
)

// StoreRepository defines the interface for store data operations.
type StoreRepository interface {
	Create(store *domain.Store) error
	FindByID(id uint) (*domain.Store, error)
	Update(store *domain.Store) error
	Delete(id uint) error
	FindAll(filter domain.StoreFilter) ([]domain.Store, int64, error)
	FindByUserID(userID uint) ([]domain.Store, error)
	// Add other store-related methods here as needed
}

// storeRepository implements the StoreRepository interface.
type storeRepository struct {
	db *gorm.DB
}

// NewStoreRepository creates a new instance of StoreRepository.
func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db: db}
}

// Create a new store in the database.
func (r *storeRepository) Create(store *domain.Store) error {
	return r.db.Create(store).Error
}

// FindByID retrieves a store by its ID.
func (r *storeRepository) FindByID(id uint) (*domain.Store, error) {
	var store domain.Store
	err := r.db.First(&store, id).Error
	return &store, err
}

// Update an existing store in the database.
func (r *storeRepository) Update(store *domain.Store) error {
	return r.db.Save(store).Error
}

// Delete a store by its ID (soft delete).
func (r *storeRepository) Delete(id uint) error {
	var store domain.Store
	return r.db.Delete(&store, id).Error
}

// FindAll retrieves all stores from the database with pagination and filtering.
func (r *storeRepository) FindAll(filter domain.StoreFilter) ([]domain.Store, int64, error) {
	var stores []domain.Store
	query := r.db.Model(&domain.Store{})

	// Apply search filter
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name LIKE ? OR description LIKE ? OR address LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Get total count before applying pagination
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.Limit
	query = query.Limit(filter.Limit).Offset(offset)

	err := query.Find(&stores).Error
	return stores, totalCount, err
}

// FindByUserID retrieves all stores for a specific user.
func (r *storeRepository) FindByUserID(userID uint) ([]domain.Store, error) {
	var stores []domain.Store
	err := r.db.Where("user_id = ?", userID).Find(&stores).Error
	return stores, err
}