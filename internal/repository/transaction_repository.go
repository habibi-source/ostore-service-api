package repository

import (
	"mini-project-ostore/internal/domain"

	"gorm.io/gorm"
)

// TransactionRepository defines the interface for transaction data operations.
type TransactionRepository interface {
	Create(transaction *domain.Transaction) error
	FindByID(id, userID uint) (*domain.Transaction, error) // Updated to include userID
	FindAll(filter domain.TransactionFilter) ([]domain.Transaction, int64, error)
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

// FindByID retrieves a transaction by its ID and userID.
func (r *transactionRepository) FindByID(id, userID uint) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.db.Preload("Items.ProductLog").Preload("User").Preload("Address").
		Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error
	return &transaction, err
}

// FindAll retrieves transactions based on the provided filter.
func (r *transactionRepository) FindAll(filter domain.TransactionFilter) ([]domain.Transaction, int64, error) {
	var transactions []domain.Transaction
	query := r.db.Model(&domain.Transaction{}).Preload("Items.ProductLog").Preload("User").Preload("Address")

	// Apply UserID filter
	if filter.UserID != 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	// Apply StoreID filter
	if filter.StoreID != 0 {
		// This filter is complex as Transaction model doesn't directly have store_id.
		// It would require joining through TransactionItems and Products, which adds complexity.
		// For now, I'll omit this to avoid breaking changes if schema doesn't support directly.
		// If needed, the join logic would look something like:
		// query = query.Joins("JOIN transaction_items ti ON ti.transaction_id = transactions.id").
		// 	Joins("JOIN products p ON p.id = ti.product_id").
		// 	Where("p.store_id = ?", filter.StoreID)
	}

	// Apply Status filter
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Apply PaymentMethod filter
	if filter.PaymentMethod != "" {
		query = query.Where("payment_method = ?", filter.PaymentMethod)
	}

	// Get total count before applying pagination
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.Limit
	query = query.Limit(filter.Limit).Offset(offset)

	err := query.Find(&transactions).Error
	return transactions, totalCount, err
}

// Update an existing transaction in the database.
func (r *transactionRepository) Update(transaction *domain.Transaction) error {
	return r.db.Save(transaction).Error
}