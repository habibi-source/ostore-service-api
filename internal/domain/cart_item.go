package domain

import (
	"time"

	"gorm.io/gorm"
)

type CartItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	ProductID uint           `gorm:"not null" json:"product_id"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	Price     float64        `gorm:"type:decimal(10,2);not null" json:"price"` // Price at the time of adding to cart
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product   Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}