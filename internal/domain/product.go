package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	StoreID     uint           `gorm:"not null" json:"store_id"`
	CategoryID  uint           `gorm:"not null" json:"category_id"`
	SKU         string         `gorm:"size:50;uniqueIndex;not null" json:"sku"`
	Slug        string         `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Name        string         `gorm:"size:200;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int            `gorm:"not null" json:"stock"`
	Weight      float64        `gorm:"type:decimal(10,2)" json:"weight"`
	Images      string         `gorm:"type:text" json:"images"` // JSON array of image URLs
	IsAvailable bool           `gorm:"default:true" json:"is_available"`
	Store       Store          `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}