package domain

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"not null" json:"user_id"`
	AddressID       uint           `gorm:"not null" json:"address_id"`
	InvoiceNumber   string         `gorm:"size:100;uniqueIndex;not null" json:"invoice_number"`
	TotalAmount     float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	ShippingCost    float64        `gorm:"type:decimal(10,2);not null" json:"shipping_cost"`
	PaymentMethod   string         `gorm:"size:50;not null" json:"payment_method"` // e.g., credit_card, bank_transfer, COD
	Status          string         `gorm:"size:50;not null" json:"status"`         // pending, paid, shipped, completed, cancelled
	ShippingCourier string         `gorm:"size:50" json:"shipping_courier"`
	ShippingTracking string        `gorm:"size:100" json:"shipping_tracking"`
	ConfirmedAt     gorm.DeletedAt `json:"confirmed_at"` // Timestamp when order is confirmed by seller
	PaidAt          gorm.DeletedAt `json:"paid_at"`      // Timestamp when payment is received
	ShippedAt       gorm.DeletedAt `json:"shipped_at"`   // Timestamp when order is shipped
	CompletedAt     gorm.DeletedAt `json:"completed_at"` // Timestamp when order is completed (received by customer)
	CancelledAt     gorm.DeletedAt `json:"cancelled_at"` // Timestamp when order is cancelled
	User            User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Address         Address        `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	Items           []TransactionItem `gorm:"foreignKey:TransactionID" json:"items"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type TransactionItem struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	TransactionID uint         `gorm:"not null" json:"transaction_id"`
	ProductID     uint         `gorm:"not null" json:"product_id"`
	Quantity      int          `gorm:"not null" json:"quantity"`
	Price         float64      `gorm:"type:decimal(10,2);not null" json:"price"`
	ProductLog    ProductLog   `gorm:"foreignKey:TransactionItemID" json:"product_log"`
}

type ProductLog struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	TransactionItemID  uint      `gorm:"not null" json:"transaction_item_id"`
	ProductName        string    `gorm:"size:200;not null" json:"product_name"`
	ProductDescription string    `gorm:"type:text" json:"product_description"`
	ProductPrice       float64   `gorm:"type:decimal(10,2);not null" json:"product_price"`
	ProductWeight      float64   `gorm:"type:decimal(10,2)" json:"product_weight"`
	ProductImages      string    `gorm:"type:text" json:"product_images"`
	CreatedAt          time.Time `json:"created_at"`
}