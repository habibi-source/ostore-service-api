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

// PaginatedTransactionResponse defines the structure for a paginated list of transactions.
type PaginatedTransactionResponse struct {
	Transactions []domain.Transaction `json:"transactions"`
	Page         int                  `json:"page"`
	Limit        int                  `json:"limit"`
	TotalCount   int64                `json:"total_count"`
	TotalPages   int                  `json:"total_pages"`
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

// GetUserTransactions retrieves all transactions for the logged-in user with pagination and filtering.
func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	// Get user_id from JWT token set by middleware
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

	var filter domain.TransactionFilter
	filter.UserID = userID // Set UserID from the authenticated user

	// Parse pagination parameters
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			filter.Page = page
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}
	filter.SetDefaults() // Apply default page and limit if not set

	// Parse filtering parameters
	filter.Status = c.Query("status")
	filter.PaymentMethod = c.Query("payment_method")

	if storeIDStr := c.Query("store_id"); storeIDStr != "" {
		if storeID, err := strconv.ParseUint(storeIDStr, 10, 32); err == nil {
			filter.StoreID = uint(storeID)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store_id format"})
			return
		}
	}

	// Retrieve all transactions for this user with filter and pagination
	transactions, totalCount, err := h.transactionUC.GetUserTransactions(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := 0
	if filter.Limit > 0 {
		totalPages = int((totalCount + int64(filter.Limit) - 1) / int64(filter.Limit))
	}

	c.JSON(http.StatusOK, domain.StandardPaginatedResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Data: PaginatedTransactionResponse{
			Transactions: transactions,
			Page:         filter.Page,
			Limit:        filter.Limit,
			TotalCount:   totalCount,
			TotalPages:   totalPages,
		},
	})
}

// GetTransaction retrieves a single transaction by its ID for the authenticated user.
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	// Get user_id from JWT token set by middleware
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

	transactionIDStr := c.Param("id")
	transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	// Pass both transactionID and userID to the use case for ownership validation
	transaction, err := h.transactionUC.GetByID(uint(transactionID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}) // This will also catch "not found" if not owned by user
		return
	}
	c.JSON(http.StatusOK, transaction)
}