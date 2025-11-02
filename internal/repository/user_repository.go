// internal/repository/user_repository.go
package repository

import (
	"mini-project-ostore/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindByPhone(phone string) (*domain.User, error)
	Update(user *domain.User) error
	EmailExists(email string, excludeID uint) (bool, error)
	PhoneExists(phone string, excludeID uint) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByPhone(phone string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	return &user, err
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) EmailExists(email string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&domain.User{}).Where("email = ?", email)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *userRepository) PhoneExists(phone string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&domain.User{}).Where("phone = ?", phone)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
