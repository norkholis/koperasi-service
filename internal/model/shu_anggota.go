package model

import (
	"time"

	"gorm.io/gorm"
)

// SHUAnggotaRecord represents the shu_anggota table for saving individual member SHU calculations
type SHUAnggotaRecord struct {
	ID          uint           `json:"id_shu_anggota" gorm:"primaryKey;column:id_shu_anggota"`
	SHUID       uint           `json:"id_shu" gorm:"column:id_shu;not null"`
	UserID      uint           `json:"id_anggota" gorm:"column:id_anggota;not null"`
	JumlahModal float64        `json:"jumlah_modal" gorm:"column:jumlah_modal;type:decimal(15,2);not null"`
	JumlahUsaha float64        `json:"jumlah_usaha" gorm:"column:jumlah_usaha;type:decimal(15,2);not null"`
	SHUDiterima float64        `json:"shu_diterima" gorm:"column:shu_diterima;type:decimal(15,2);not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	SHU  SHUTahunan `json:"shu,omitempty" gorm:"foreignKey:SHUID;references:ID"`
	User User       `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

// TableName specifies the table name for SHUAnggotaRecord
func (SHUAnggotaRecord) TableName() string {
	return "shu_anggota"
}
