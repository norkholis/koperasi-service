package model

import (
	"time"

	"gorm.io/gorm"
)

// SHUTahunan represents annual profit sharing (Sisa Hasil Usaha) record
type SHUTahunan struct {
	gorm.Model
	Tahun         int       `gorm:"not null;index" json:"tahun"`
	TotalSHU      float64   `gorm:"type:decimal(15,2);not null" json:"total_shu"`
	TanggalHitung time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"tanggal_hitung"`
	Status        string    `gorm:"type:varchar(20);check:status IN ('draft', 'final')" json:"status"`
}

// SHUAnggota represents individual member's SHU calculation result
type SHUAnggota struct {
	UserID          uint    `json:"user_id"`
	Email           string  `json:"email"`
	TotalSimpanan   float64 `json:"total_simpanan"`
	TotalPenjualan  float64 `json:"total_penjualan"`
	JasaModal       float64 `json:"jasa_modal"`
	JasaUsaha       float64 `json:"jasa_usaha"`
	TotalSHUAnggota float64 `json:"total_shu_anggota"`
}

// SHUReport represents the complete SHU calculation report
type SHUReport struct {
	Tahun             int          `json:"tahun"`
	TotalSHUKoperasi  float64      `json:"total_shu_koperasi"`
	PersenJasaModal   float64      `json:"persen_jasa_modal"`
	PersenJasaUsaha   float64      `json:"persen_jasa_usaha"`
	TotalSimpananAll  float64      `json:"total_simpanan_all"`
	TotalPenjualanAll float64      `json:"total_penjualan_all"`
	TanggalHitung     time.Time    `json:"tanggal_hitung"`
	DetailAnggota     []SHUAnggota `json:"detail_anggota"`
}
