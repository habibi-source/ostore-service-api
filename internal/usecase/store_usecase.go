package usecase

import (
	"errors"
	"mini-project-ostore/internal/domain"
	"mini-project-ostore/internal/repository"
)

// StoreUseCase defines the interface for store-related business logic.
type StoreUseCase interface {
	Create(store *domain.Store) error
	GetByID(id uint) (*domain.Store, error)
	Update(store *domain.Store) error
	Delete(id uint) error
	GetByUserID(userID uint) ([]domain.Store, error)
	GetAll() ([]domain.Store, error)
}

// storeUseCase implements the StoreUseCase interface.
type storeUseCase struct {
	storeRepo repository.StoreRepository
	userRepo  repository.UserRepository // For validating UserID if needed
}

// NewStoreUseCase creates a new instance of StoreUseCase.
func NewStoreUseCase(storeRepo repository.StoreRepository, userRepo repository.UserRepository) StoreUseCase {
	return &storeUseCase{storeRepo: storeRepo, userRepo: userRepo}
}

// Create a new store.
func (uc *storeUseCase) Create(store *domain.Store) error {
	// Optional: Check if the user exists before creating a store for them.
	// This ensures referential integrity at the application level.
	_, err := uc.userRepo.FindByID(store.UserID)
	if err != nil {
		return errors.New("user not found for the given UserID")
	}
	return uc.storeRepo.Create(store)
}

// GetByID retrieves a store by its ID.
func (uc *storeUseCase) GetByID(id uint) (*domain.Store, error) {
	return uc.storeRepo.FindByID(id)
}

// Update an existing store.
func (uc *storeUseCase) Update(store *domain.Store) error {
	// First, check if the store exists.
	existingStore, err := uc.storeRepo.FindByID(store.ID)
	if err != nil {
		return errors.New("store not found")
	}

	// Update only the mutable fields if they are provided in the update request.
	// This prevents overwriting existing data with empty strings if a field is not sent.
	if store.Name != "" {
		existingStore.Name = store.Name
	}
	if store.Description != "" {
		existingStore.Description = store.Description
	}
	if store.Address != "" {
		existingStore.Address = store.Address
	}
	if store.Phone != "" {
		existingStore.Phone = store.Phone
	}
	// Update PhotoProfile only if a new photo is provided
	if store.PhotoProfile != "" {
		existingStore.PhotoProfile = store.PhotoProfile
	}

	return uc.storeRepo.Update(existingStore)
}

// Delete a store by its ID.
func (uc *storeUseCase) Delete(id uint) error {
	// First, check if the store exists.
	_, err := uc.storeRepo.FindByID(id)
	if err != nil {
		return errors.New("store not found")
	}
	return uc.storeRepo.Delete(id)
}

// GetByUserID retrieves all stores owned by a specific user.
func (uc *storeUseCase) GetByUserID(userID uint) ([]domain.Store, error) {
	// Check if the user exists
	_, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return uc.storeRepo.FindByUserID(userID)
}

// GetAll retrieves all stores.
func (uc *storeUseCase) GetAll() ([]domain.Store, error) {
	return uc.storeRepo.FindAll()
}