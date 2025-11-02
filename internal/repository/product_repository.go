package repository

import (
	"mini-project-ostore/internal/domain"

	"gorm.io/gorm"
)

// ProductRepository defines the interface for product data operations.
type ProductRepository interface {
	Create(product *domain.Product) error
	FindByID(id uint) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uint) error
	GetProducts(filter domain.ProductFilter) ([]domain.Product, int64, error)
	SKUExists(sku string, excludeID uint) (bool, error)
	SlugExists(slug string, excludeID uint) (bool, error)
}

// productRepository implements the ProductRepository interface.
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new instance of ProductRepository.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create a new product in the database.
func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

// FindByID retrieves a product by its ID.
func (r *productRepository) FindByID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.First(&product, id).Error
	return &product, err
}

// Update an existing product in the database.
func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

// Delete a product by its ID (soft delete).
func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

// GetProducts retrieves products based on the provided filter.
func (r *productRepository) GetProducts(filter domain.ProductFilter) ([]domain.Product, int64, error) {
	var products []domain.Product
	query := r.db.Model(&domain.Product{})

	// Apply search filter
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name LIKE ? OR description LIKE ? OR sku LIKE ? OR slug LIKE ?", searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Apply category ID filter
	if filter.CategoryID != 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}

	// Apply price range filter
	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}

	// Apply store ID filter
	if filter.StoreID != 0 {
		query = query.Where("products.store_id = ?", filter.StoreID) // Use products.store_id to avoid ambiguity if JOINs are present
	}

	// Apply user ID filter (via stores)
	if filter.UserID != 0 {
		query = query.Joins("JOIN stores ON products.store_id = stores.id").Where("stores.user_id = ?", filter.UserID)
	}

	// Get total count before applying pagination
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.Limit
	query = query.Limit(filter.Limit).Offset(offset)

	err := query.Find(&products).Error
	return products, totalCount, err
}



// SKUExists checks if a product with the given SKU already exists.
func (r *productRepository) SKUExists(sku string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&domain.Product{}).Where("sku = ?", sku)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// SlugExists checks if a product with the given slug already exists.
func (r *productRepository) SlugExists(slug string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&domain.Product{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}