package domain

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProductID uint           `gorm:"not null" json:"product_id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Rating    int            `gorm:"type:int;not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment   string         `gorm:"type:text" json:"comment"`
	Product   Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}