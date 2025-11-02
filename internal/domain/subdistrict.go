package domain

import (
	"time"

	"gorm.io/gorm"
)

type Subdistrict struct {
	ID        string         `gorm:"primaryKey;size:100" json:"id"` // Use string ID for region codes
	CityID    string         `gorm:"size:100;not null" json:"city_id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	City      City           `gorm:"foreignKey:CityID" json:"city,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}