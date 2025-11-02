package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Email       string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Phone       string         `gorm:"size:20;uniqueIndex" json:"phone"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	IsAdmin     bool           `gorm:"default:false" json:"is_admin"`
	Stores      []Store        `gorm:"foreignKey:UserID" json:"stores,omitempty"`
	Addresses   []Address      `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}