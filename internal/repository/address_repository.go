package repository

import (
	"mini-project-ostore/internal/domain"

	"gorm.io/gorm"
)

// AddressRepository defines the interface for address data operations.
type AddressRepository interface {
	Create(address *domain.Address) error
	FindByID(id uint) (*domain.Address, error)
	Update(address *domain.Address) error
	Delete(id uint) error
	GetUserAddresses(userID uint) ([]domain.Address, error)
	// Add other address-related methods here as needed
}

// addressRepository implements the AddressRepository interface.
type addressRepository struct {
	db *gorm.DB
}

// NewAddressRepository creates a new instance of AddressRepository.
func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

// Create a new address in the database.
func (r *addressRepository) Create(address *domain.Address) error {
	return r.db.Create(address).Error
}

// FindByID retrieves an address by its ID.
func (r *addressRepository) FindByID(id uint) (*domain.Address, error) {
	var address domain.Address
	err := r.db.First(&address, id).Error
	return &address, err
}

// Update an existing address in the database.
func (r *addressRepository) Update(address *domain.Address) error {
	return r.db.Save(address).Error
}

// Delete an address by its ID (soft delete).
func (r *addressRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Address{}, id).Error
}

// GetUserAddresses retrieves all addresses for a specific user.
func (r *addressRepository) GetUserAddresses(userID uint) ([]domain.Address, error) {
	var addresses []domain.Address
	err := r.db.Where("user_id = ?", userID).Find(&addresses).Error
	return addresses, err
}