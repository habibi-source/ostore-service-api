package domain

import (
	"time"

	"gorm.io/gorm"
)

type City struct {
	ID          string         `gorm:"primaryKey;size:100" json:"id"` // Use string ID for region codes
	ProvinceID  string         `gorm:"size:100;not null" json:"province_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Province    Province       `gorm:"foreignKey:ProvinceID" json:"province,omitempty"`
	Subdistricts []Subdistrict `gorm:"foreignKey:CityID" json:"subdistricts,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
