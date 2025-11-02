package usecase

import (
	"errors"
	"mini-project-ostore/internal/domain"
	"mini-project-ostore/internal/repository"

	"github.com/google/uuid"
	"time"
)

type TransactionUseCase interface {
	Create(transaction *domain.Transaction) error
	GetByID(id, userID uint) (*domain.Transaction, error) // Updated to include userID
	GetUserTransactions(filter domain.TransactionFilter) ([]domain.Transaction, int64, error)
}

type transactionUseCase struct {
	transactionRepo repository.TransactionRepository
	productRepo     repository.ProductRepository
	userRepo        repository.UserRepository
	addressRepo     repository.AddressRepository
}

func NewTransactionUseCase(transactionRepo repository.TransactionRepository, productRepo repository.ProductRepository, userRepo repository.UserRepository, addressRepo repository.AddressRepository) TransactionUseCase {
	return &transactionUseCase{transactionRepo: transactionRepo, productRepo: productRepo, userRepo: userRepo, addressRepo: addressRepo}
}

func (uc *transactionUseCase) Create(transaction *domain.Transaction) error {
	_, err := uc.userRepo.FindByID(transaction.UserID)
	if err != nil {
		return errors.New("user not found for the given UserID")
	}

	_, err = uc.addressRepo.FindByID(transaction.AddressID)
	if err != nil {
		return errors.New("address not found for the given AddressID")
	}

	var totalAmount float64

	for i, item := range transaction.Items {
		product, err := uc.productRepo.FindByID(item.ProductID)
		if err != nil {
			return errors.New("product not found for item")
		}
		if product.Stock < item.Quantity {
			return errors.New("not enough stock for product")
		}
		product.Stock -= item.Quantity
		err = uc.productRepo.Update(product)
		if err != nil {
			return errors.New("failed to update product stock")
		}

		// Calculate total amount
		totalAmount += item.Price * float64(item.Quantity)

		// Populate ProductLog
		transaction.Items[i].ProductLog = domain.ProductLog{
			ProductName:        product.Name,
			ProductDescription: product.Description,
			ProductPrice:       product.Price,
			ProductWeight:      product.Weight,
			ProductImages:      product.Images,
			CreatedAt:          time.Now(),
		}
	}

	transaction.TotalAmount = totalAmount + transaction.ShippingCost
	transaction.InvoiceNumber = uuid.New().String()
	transaction.Status = "pending"
	// PaymentMethod and ShippingCourier/ShippingTracking are expected to be set by the handler.
	// ConfirmedAt, PaidAt, ShippedAt, CompletedAt, CancelledAt are nil by default and updated later by status changes.

	return uc.transactionRepo.Create(transaction)
}

// GetByID retrieves a transaction by its ID and userID for ownership validation.
func (uc *transactionUseCase) GetByID(id, userID uint) (*domain.Transaction, error) {
	// Optionally, check if the user exists before querying the transaction
	_, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found or an error occurred while checking user existence")
	}
	return uc.transactionRepo.FindByID(id, userID)
}

// GetUserTransactions retrieves transactions for a specific user with pagination and filtering.
func (uc *transactionUseCase) GetUserTransactions(filter domain.TransactionFilter) ([]domain.Transaction, int64, error) {
	// Validate UserID from filter
	_, err := uc.userRepo.FindByID(filter.UserID)
	if err != nil {
		return nil, 0, errors.New("user not found or an error occurred while checking user existence")
	}
	return uc.transactionRepo.FindAll(filter)
}