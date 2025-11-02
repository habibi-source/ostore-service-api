package usecase

import (
	"errors"
	"mini-project-ostore/internal/domain"
	"mini-project-ostore/internal/repository"
)

// AddressUseCase defines the interface for address-related business logic.
type AddressUseCase interface {
	Create(address *domain.Address) error
	GetByID(id uint) (*domain.Address, error)
	Update(address *domain.Address) error
	Delete(id uint) error
	GetByUserID(userID uint) ([]domain.Address, error)
}

// addressUseCase implements the AddressUseCase interface.
type addressUseCase struct {
	addressRepo repository.AddressRepository
	userRepo    repository.UserRepository // For validating UserID
}

// NewAddressUseCase creates a new instance of AddressUseCase.
func NewAddressUseCase(addressRepo repository.AddressRepository, userRepo repository.UserRepository) AddressUseCase {
	return &addressUseCase{addressRepo: addressRepo, userRepo: userRepo}
}

// Create a new address.
func (uc *addressUseCase) Create(address *domain.Address) error {
	// Optional: Check if the user exists before creating an address for them.
	_, err := uc.userRepo.FindByID(address.UserID)
	if err != nil {
		return errors.New("user not found for the given UserID")
	}
	return uc.addressRepo.Create(address)
}

// GetByID retrieves an address by its ID.
func (uc *addressUseCase) GetByID(id uint) (*domain.Address, error) {
	return uc.addressRepo.FindByID(id)
}

// Update an existing address.
func (uc *addressUseCase) Update(address *domain.Address) error {
	// First, check if the address exists.
	existingAddress, err := uc.addressRepo.FindByID(address.ID)
	if err != nil {
		return errors.New("address not found")
	}

	// Ensure the UserID is not changed during update, or handle it explicitly if allowed.
	if existingAddress.UserID != address.UserID {
		return errors.New("cannot change UserID of an existing address")
	}

	// Update only the mutable fields.
	existingAddress.Label = address.Label
	existingAddress.ReceiverName = address.ReceiverName
	existingAddress.Phone = address.Phone
	existingAddress.ProvinceID = address.ProvinceID
	existingAddress.CityID = address.CityID
	existingAddress.SubDistrictID = address.SubDistrictID
	existingAddress.Detail = address.Detail
	existingAddress.PostalCode = address.PostalCode
	existingAddress.IsPrimary = address.IsPrimary

	return uc.addressRepo.Update(existingAddress)
}

// Delete an address by its ID.
func (uc *addressUseCase) Delete(id uint) error {
	// First, check if the address exists.
	_, err := uc.addressRepo.FindByID(id)
	if err != nil {
		return errors.New("address not found")
	}
	return uc.addressRepo.Delete(id)
}

// GetByUserID retrieves all addresses for a specific user.
func (uc *addressUseCase) GetByUserID(userID uint) ([]domain.Address, error) {
	// Check if the user exists
	_, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return uc.addressRepo.GetUserAddresses(userID)
}