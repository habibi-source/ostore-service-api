package domain

import (
	"time"

	"gorm.io/gorm"
)

type Province struct {
	ID        string         `gorm:"primaryKey;size:100" json:"id"` // Use string ID for region codes
	Name      string         `gorm:"size:100;not null" json:"name"`
	Cities    []City         `gorm:"foreignKey:ProvinceID" json:"cities,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}