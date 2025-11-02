package usecase

import (
	"errors"
	"mini-project-ostore/internal/domain"
	"mini-project-ostore/internal/repository"
)

// ProductUseCase defines the interface for product-related business logic.
type ProductUseCase interface {
	Create(product *domain.Product, userID uint) error
	GetByID(id uint) (*domain.Product, error)
	Update(product *domain.Product, userID uint) error
	Delete(id uint, userID uint) error
	GetProducts(filter domain.ProductFilter) ([]domain.Product, int64, error)
	GetUserProducts(filter domain.ProductFilter) ([]domain.Product, int64, error)
}

// productUseCase implements the ProductUseCase interface.
type productUseCase struct {
	productRepo  repository.ProductRepository
	storeRepo    repository.StoreRepository
	userRepo     repository.UserRepository
	categoryRepo repository.CategoryRepository
}

// NewProductUseCase creates a new instance of ProductUseCase.
func NewProductUseCase(
	productRepo repository.ProductRepository,
	storeRepo repository.StoreRepository,
	userRepo repository.UserRepository,
	categoryRepo repository.CategoryRepository,
) ProductUseCase {
	return &productUseCase{
		productRepo:  productRepo,
		storeRepo:    storeRepo,
		userRepo:     userRepo,
		categoryRepo: categoryRepo,
	}
}

// Create creates a new product for a user's store.
func (uc *productUseCase) Create(product *domain.Product, userID uint) error {
	store, err := uc.storeRepo.FindByID(product.StoreID)
	if err != nil {
		return errors.New("store not found for the given StoreID")
	}

	// ✅ Validasi: store harus milik user yang login
	if store.UserID != userID {
		return errors.New("unauthorized: cannot create product for another user's store")
	}

	// Check if SKU already exists
	skuExists, err := uc.productRepo.SKUExists(product.SKU, 0)
	if err != nil {
		return err
	}
	if skuExists {
		return errors.New("product SKU already exists")
	}

	// Check if Slug already exists
	slugExists, err := uc.productRepo.SlugExists(product.Slug, 0)
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
func (uc *productUseCase) Update(product *domain.Product, userID uint) error {
	existingProduct, err := uc.productRepo.FindByID(product.ID)
	if err != nil {
		return errors.New("product not found")
	}

	store, err := uc.storeRepo.FindByID(existingProduct.StoreID)
	if err != nil {
		return errors.New("store not found for product")
	}

	// ✅ Validasi: user hanya boleh ubah produk dari toko miliknya
	if store.UserID != userID {
		return errors.New("unauthorized: cannot update product from another user's store")
	}

	// Field update
	if product.StoreID != 0 && product.StoreID != existingProduct.StoreID {
		newStore, err := uc.storeRepo.FindByID(product.StoreID)
		if err != nil {
			return errors.New("new StoreID not found")
		}
		if newStore.UserID != userID {
			return errors.New("unauthorized: cannot move product to another user's store")
		}
		existingProduct.StoreID = product.StoreID
	}

	if product.CategoryID != 0 && product.CategoryID != existingProduct.CategoryID {
		if _, err := uc.categoryRepo.FindByID(product.CategoryID); err != nil {
			return errors.New("new CategoryID not found")
		}
		existingProduct.CategoryID = product.CategoryID
	}

	if product.SKU != "" && product.SKU != existingProduct.SKU {
		exists, _ := uc.productRepo.SKUExists(product.SKU, product.ID)
		if exists {
			return errors.New("product SKU already exists")
		}
		existingProduct.SKU = product.SKU
	}

	if product.Slug != "" && product.Slug != existingProduct.Slug {
		exists, _ := uc.productRepo.SlugExists(product.Slug, product.ID)
		if exists {
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

	return uc.productRepo.Update(existingProduct)
}

// Delete removes a product if it belongs to the user's store.
func (uc *productUseCase) Delete(id uint, userID uint) error {
	product, err := uc.productRepo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	store, err := uc.storeRepo.FindByID(product.StoreID)
	if err != nil {
		return errors.New("store not found for product")
	}

	// ✅ Validasi: hanya owner store bisa menghapus produk
	if store.UserID != userID {
		return errors.New("unauthorized: cannot delete another user's product")
	}

	return uc.productRepo.Delete(id)
}

// GetProducts retrieves products based on filters (public endpoint).
func (uc *productUseCase) GetProducts(filter domain.ProductFilter) ([]domain.Product, int64, error) {
	return uc.productRepo.GetProducts(filter)
}

// GetUserProducts retrieves products for stores owned by a specific user.
func (uc *productUseCase) GetUserProducts(filter domain.ProductFilter) ([]domain.Product, int64, error) {
	if _, err := uc.userRepo.FindByID(filter.UserID); err != nil {
		return nil, 0, errors.New("user not found")
	}
	return uc.productRepo.GetProducts(filter)
}
