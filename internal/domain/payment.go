package domain

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TransactionID uint           `gorm:"uniqueIndex;not null" json:"transaction_id"` // Foreign key to Transaction
	Amount        float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	PaymentMethod string         `gorm:"size:50;not null" json:"payment_method"` // e.g., credit_card, bank_transfer, COD
	Status        string         `gorm:"size:50;not null" json:"status"`         // pending, completed, failed, refunded
	PaymentDate   time.Time      `gorm:"not null" json:"payment_date"`
	Transaction   Transaction    `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}