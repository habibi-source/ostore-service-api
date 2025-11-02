package usecase

import (
	"errors"
	"mini-project-ostore/internal/domain"
	"mini-project-ostore/internal/repository"
)

// ProductUseCase defines the interface for product-related business logic.
type ProductUseCase interface {
	Create(product *domain.Product) error
	GetByID(id uint) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uint) error
	GetAll() ([]domain.Product, error)
	GetUserProducts(userID uint) ([]domain.Product, error)
}

// productUseCase implements the ProductUseCase interface.
type productUseCase struct {
	productRepo  repository.ProductRepository
	storeRepo    repository.StoreRepository    // For validating StoreID
	userRepo     repository.UserRepository     // For validating UserID in GetUserProducts
	categoryRepo repository.CategoryRepository // For validating CategoryID
}

// NewProductUseCase creates a new instance of ProductUseCase.
func NewProductUseCase(productRepo repository.ProductRepository, storeRepo repository.StoreRepository, userRepo repository.UserRepository, categoryRepo repository.CategoryRepository) ProductUseCase {
	return &productUseCase{productRepo: productRepo, storeRepo: storeRepo, userRepo: userRepo, categoryRepo: categoryRepo}
}

// Create a new product.
func (uc *productUseCase) Create(product *domain.Product) error {
	// Check if the store exists
	_, err := uc.storeRepo.FindByID(product.StoreID)
	if err != nil {
		return errors.New("store not found for the given StoreID")
	}

	// Check if SKU already exists
	skuExists, err := uc.productRepo.SKUExists(product.SKU, 0) // 0 for no exclude ID
	if err != nil {
		return err
	}
	if skuExists {
		return errors.New("product SKU already exists")
	}

	// Check if Slug already exists
	slugExists, err := uc.productRepo.SlugExists(product.Slug, 0) // 0 for no exclude ID
	if err != nil {
		return err
	}
	if slugExists {
		return errors.New("product slug already exists")
	}

	return uc.productRepo.Create(product)
}

// GetByID retrieves a product by its ID.
func (uc *productUseCase) GetByID(id uint) (*domain.Product, error) {
	return uc.productRepo.FindByID(id)
}

// Update an existing product.
func (uc *productUseCase) Update(product *domain.Product) error {
	// First, check if the product exists.
	existingProduct, err := uc.productRepo.FindByID(product.ID)
	if err != nil {
		return errors.New("product not found")
	}

	// Update fields only if they are explicitly provided in the input 'product'
	// Check and update StoreID if provided and different
	if product.StoreID != 0 && existingProduct.StoreID != product.StoreID {
		// Validate if the new store exists
		_, err := uc.storeRepo.FindByID(product.StoreID)
		if err != nil {
			return errors.New("new StoreID not found")
		}
		existingProduct.StoreID = product.StoreID
	}

	// Check and update CategoryID if provided and different
	if product.CategoryID != 0 && existingProduct.CategoryID != product.CategoryID {
		// Validate if the new category exists
		_, err := uc.categoryRepo.FindByID(product.CategoryID)
		if err != nil {
			return errors.New("new CategoryID not found")
		}
		existingProduct.CategoryID = product.CategoryID
	}

	// Check and update SKU if provided and different
	if product.SKU != "" && existingProduct.SKU != product.SKU {
		// Check if new SKU already exists for another product
		skuExists, err := uc.productRepo.SKUExists(product.SKU, product.ID)
		if err != nil {
			return err
		}
		if skuExists {
			return errors.New("product SKU already exists")
		}
		existingProduct.SKU = product.SKU
	}

	// Check and update Slug if provided and different
	if product.Slug != "" && existingProduct.Slug != product.Slug {
		// Check if new Slug already exists for another product
		slugExists, err := uc.productRepo.SlugExists(product.Slug, product.ID)
		if err != nil {
			return err
		}
		if slugExists {
			return errors.New("product slug already exists")
		}
		existingProduct.Slug = product.Slug
	}

	if product.Name != "" {
		existingProduct.Name = product.Name
	}
	if product.Description != "" {
		existingProduct.Description = product.Description
	}
	if product.Price != 0 {
		existingProduct.Price = product.Price
	}
	if product.Stock != 0 {
		existingProduct.Stock = product.Stock
	}
	if product.Weight != 0 {
		existingProduct.Weight = product.Weight
	}
	if product.Images != "" {
		existingProduct.Images = product.Images
	}
	// IsAvailable is not directly updatable via product_handler.go's UpdateProductRequest.
	// If it needs to be updated, it should be added to UpdateProductRequest.

	return uc.productRepo.Update(existingProduct)
}

// Delete a product by its ID.
func (uc *productUseCase) Delete(id uint) error {
	// First, check if the product exists.
	_, err := uc.productRepo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}
	return uc.productRepo.Delete(id)
}

// GetAll retrieves all products.
func (uc *productUseCase) GetAll() ([]domain.Product, error) {
	return uc.productRepo.GetAll()
}

// GetUserProducts retrieves all products for stores owned by a specific user.
func (uc *productUseCase) GetUserProducts(userID uint) ([]domain.Product, error) {
	// Check if the user exists
	_, err := uc.userRepo.FindByID(userID) // This line will cause an error because userRepo is not defined in productUseCase
	if err != nil {
		return nil, errors.New("user not found")
	}
	return uc.productRepo.GetUserProducts(userID)
}