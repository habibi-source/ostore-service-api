package domain

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null" json:"user_id"`
	Label         string         `gorm:"size:50;not null" json:"label"`
	ReceiverName  string         `gorm:"size:100;not null" json:"receiver_name"`
	Phone         string         `gorm:"size:20;not null" json:"phone"`
	ProvinceID    uint         `gorm:"size:100;not null" json:"province_id"`    // Reference to Province table
	CityID        uint         `gorm:"size:100;not null" json:"city_id"`        // Reference to City table
	SubDistrictID uint         `gorm:"size:100;not null" json:"subdistrict_id"` // Reference to Subdistrict table
	Detail        string         `gorm:"type:text;not null" json:"detail"`        // Detailed address
	PostalCode    string         `gorm:"size:10;not null" json:"postal_code"`
	IsPrimary     bool           `gorm:"default:false" json:"is_primary"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
