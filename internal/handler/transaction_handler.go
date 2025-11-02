package handler

import (
	"net/http"
	"strconv"

	"mini-project-ostore/internal/domain"
	"mini-project-ostore/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionUC usecase.TransactionUseCase
}

func NewTransactionHandler(transactionUC usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{transactionUC: transactionUC}
}

type CreateTransactionItemRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	Price     float64 `json:"price" binding:"required,min=0"`
}

type CreateTransactionRequest struct {
	UserID          uint                           `json:"user_id" binding:"required"`
	AddressID       uint                           `json:"address_id" binding:"required"`
	ShippingCost    float64                        `json:"shipping_cost" binding:"required,min=0"`
	PaymentMethod   string                         `json:"payment_method" binding:"required"`
	ShippingCourier string                         `json:"shipping_courier"`
	Items           []CreateTransactionItemRequest `json:"items" binding:"required,min=1"`
}

// CreateTransaction handles the creation of a new transaction.
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction := &domain.Transaction{
		UserID:          req.UserID,
		AddressID:       req.AddressID,
		TotalAmount:     0, // Will be calculated in usecase
		ShippingCost:    req.ShippingCost,
		PaymentMethod:   req.PaymentMethod,
		ShippingCourier: req.ShippingCourier,
		Status:          "pending", // Default status
	}

	for _, itemReq := range req.Items {
		transaction.Items = append(transaction.Items, domain.TransactionItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			Price:     itemReq.Price,
		})
	}

	if err := h.transactionUC.Create(transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully", "transaction_id": transaction.ID, "invoice_number": transaction.InvoiceNumber})
}

// GetUserTransactions retrieves all transactions for a specific user.
// GetUserTransactions retrieves all transactions for the logged-in user.
func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	// Ambil user_id dari token JWT yang diset oleh middleware
	uid, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := uid.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Ambil semua transaksi milik user ini
	transactions, err := h.transactionUC.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Succeed to GET data",
		"data":    transactions,
	})
}


// GetTransaction retrieves a single transaction by its ID.
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	transactionIDStr := c.Param("id")
	transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := h.transactionUC.GetByID(uint(transactionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transaction)
}
