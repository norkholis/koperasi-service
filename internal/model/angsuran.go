package model

import (
	"time"

	"gorm.io/gorm"
)

// Angsuran represents an installment payment record in the system
type Angsuran struct {
	gorm.Model
	PinjamanID   uint      `gorm:"not null;index" json:"pinjaman_id"` // References pinjaman table
	AngsuranKe   int       `gorm:"not null" json:"angsuran_ke"`
	TanggalBayar time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"tanggal_bayar"`
	Pokok        float64   `gorm:"type:decimal(15,2);not null" json:"pokok"`
	Bunga        float64   `gorm:"type:decimal(15,2);not null" json:"bunga"`
	Denda        float64   `gorm:"type:decimal(15,2);default:0" json:"denda"`
	TotalBayar   float64   `gorm:"type:decimal(15,2);not null" json:"total_bayar"`
	UserID       uint      `gorm:"not null" json:"user_id"` // References users table
	Status       string    `gorm:"type:varchar(20);check:status IN ('proses', 'verified', 'kurang', 'lebih')" json:"status"`
	Pinjaman     Pinjaman  `gorm:"foreignKey:PinjamanID" json:"pinjaman,omitempty"`
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}