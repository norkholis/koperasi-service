package model

import "gorm.io/gorm"

// BungaOption represents the admin-configurable interest rate options
type BungaOption struct {
	gorm.Model
	Nama          string  `gorm:"type:varchar(50);not null" json:"nama"`    // e.g., "Bunga Rendah", "Bunga Standar"
	Persen        float64 `gorm:"type:decimal(5,2);not null" json:"persen"` // e.g., 1.00, 2.00, 2.50
	Deskripsi     string  `gorm:"type:text" json:"deskripsi"`               // Optional description
	IsActive      bool    `gorm:"default:true" json:"is_active"`            // To enable/disable options
	CreatedBy     uint    `gorm:"not null" json:"created_by"`               // Admin who created this option
	CreatedByUser User    `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
}

// TableName specifies the table name for BungaOption model
func (BungaOption) TableName() string {
	return "bunga_options"
}
