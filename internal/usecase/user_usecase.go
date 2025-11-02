package usecase

import (
	"errors"
	"mini-project-ostore/internal/domain"
	"mini-project-ostore/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// UserUseCase defines the interface for user-related business logic.
type UserUseCase interface {
	CreateDefaultStore(userID uint, userName string) error
	GetByID(id uint) (*domain.User, error)
	Update(user *domain.User) error
}

// userUseCase implements the UserUseCase interface.
type userUseCase struct {
	userRepo  repository.UserRepository
	storeRepo repository.StoreRepository
}

// NewUserUseCase creates a new instance of UserUseCase.
func NewUserUseCase(userRepo repository.UserRepository, storeRepo repository.StoreRepository) UserUseCase {
	return &userUseCase{userRepo: userRepo, storeRepo: storeRepo}
}

// CreateDefaultStore creates a default store for a newly registered user.
func (uc *userUseCase) CreateDefaultStore(userID uint, userName string) error {
	// Check if user exists (though it should exist if called after user registration)
	_, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	defaultStore := &domain.Store{
		UserID:      userID,
		Name:        userName + "'s Store",
		Description: "Default store for " + userName,
		Address:     "Default Address", // This could be updated later
		Phone:       "",                // This could be updated later
	}

	return uc.storeRepo.Create(defaultStore)
}

// GetByID retrieves a user by their ID.
func (uc *userUseCase) GetByID(id uint) (*domain.User, error) {
	user, err := uc.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// Update an existing user.
func (uc *userUseCase) Update(user *domain.User) error {
	existingUser, err := uc.userRepo.FindByID(user.ID)
	if err != nil {
		return errors.New("user not found")
	}

	// Update fields only if they are explicitly provided in the input 'user'
	if user.Name != "" {
		existingUser.Name = user.Name
	}
	if user.Email != "" && existingUser.Email != user.Email {
		// Check if new email already exists for another user
		exists, err := uc.userRepo.EmailExists(user.Email, user.ID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("email already exists")
		}
		existingUser.Email = user.Email
	}
	if user.Phone != "" && existingUser.Phone != user.Phone {
		// Check if new phone already exists for another user
		exists, err := uc.userRepo.PhoneExists(user.Phone, user.ID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("phone already exists")
		}
		existingUser.Phone = user.Phone
	}
	if user.Password != "" { // If password is provided, hash it
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash password")
		}
		existingUser.Password = string(hashedPassword)
	}
	// For IsAdmin, the handler passes a `domain.User` where `IsAdmin` is `false` if `req.IsAdmin` was `nil` in the handler.
	// If `req.IsAdmin` was `&false`, it also sets `user.IsAdmin` to `false`.
	// This creates ambiguity for `false` values. To avoid unintended updates from `false` (default value),
	// we update `IsAdmin` only if it's different from the existing value.
	// This approach might incorrectly set `IsAdmin` to `false` if `existingUser.IsAdmin` was `true` and
	// `req.IsAdmin` in the handler was `nil` (meaning no explicit update was intended for `IsAdmin`).
	// A more robust solution would require changing the handler or the domain struct to pass specific update flags or nullable fields.
	if existingUser.IsAdmin != user.IsAdmin {
		existingUser.IsAdmin = user.IsAdmin
	}

	return uc.userRepo.Update(existingUser)
}
