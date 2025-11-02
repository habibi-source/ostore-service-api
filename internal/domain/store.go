package domain

import (
	"time"

	"gorm.io/gorm"
)

type Store struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"not null" json:"user_id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	Address      string         `gorm:"type:text" json:"address"`
	Phone        string         `gorm:"size:20" json:"phone"`
	PhotoProfile string         `gorm:"size:255" json:"photo_profile"` // New field for store profile photo
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Products     []Product      `gorm:"foreignKey:StoreID" json:"products,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}