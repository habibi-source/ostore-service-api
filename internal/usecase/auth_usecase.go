// internal/usecase/auth_usecase.go
package usecase

import (
	"errors"
	"mini-project-ostore/internal/domain"
	"mini-project-ostore/internal/repository"
	"mini-project-ostore/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(user *domain.User) error
	Login(email, password string) (string, *domain.User, error)
}

type authUseCase struct {
	userRepo repository.UserRepository
}

func NewAuthUseCase(userRepo repository.UserRepository) AuthUseCase {
	return &authUseCase{userRepo: userRepo}
}

func (uc *authUseCase) Register(user *domain.User) error {
	// Check if email already exists
	exists, err := uc.userRepo.EmailExists(user.Email, 0)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	// Check if phone already exists
	if user.Phone != "" {
		exists, err = uc.userRepo.PhoneExists(user.Phone, 0)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("phone number already exists")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return uc.userRepo.Create(user)
}

func (uc *authUseCase) Login(email, password string) (string, *domain.User, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
