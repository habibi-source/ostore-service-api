package repository

import (
	"mini-project-ostore/internal/domain"

	"gorm.io/gorm"
)

// TransactionRepository defines the interface for transaction data operations.
type TransactionRepository interface {
	Create(transaction *domain.Transaction) error
	FindByID(id uint) (*domain.Transaction, error)
	FindByUserID(userID uint) ([]domain.Transaction, error)
	Update(transaction *domain.Transaction) error
	// Add other transaction-related methods here as needed
}

// transactionRepository implements the TransactionRepository interface.
type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new instance of TransactionRepository.
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// Create a new transaction in the database.
func (r *transactionRepository) Create(transaction *domain.Transaction) error {
	return r.db.Create(transaction).Error
}

// FindByID retrieves a transaction by its ID.
func (r *transactionRepository) FindByID(id uint) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.db.Preload("Items.ProductLog").Preload("User").Preload("Address").First(&transaction, id).Error
	return &transaction, err
}

// FindByUserID retrieves all transactions for a specific user.
func (r *transactionRepository) FindByUserID(userID uint) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.Where("user_id = ?", userID).Preload("Items.ProductLog").Preload("User").Preload("Address").Find(&transactions).Error
	return transactions, err
}

// Update an existing transaction in the database.
func (r *transactionRepository) Update(transaction *domain.Transaction) error {
	return r.db.Save(transaction).Error
}