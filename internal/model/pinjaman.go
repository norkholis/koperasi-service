package model

import (
	"time"

	"gorm.io/gorm"
)

// Pinjaman represents a loan record in the system with installment tracking
type Pinjaman struct {
	gorm.Model
	KodePinjaman        string       `gorm:"type:varchar(20);uniqueIndex;not null" json:"kode_pinjaman"`
	UserID              uint         `gorm:"not null" json:"user_id"` // References users table
	TanggalPinjam       time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"tanggal_pinjam"`
	JumlahPinjaman      float64      `gorm:"type:decimal(15,2);not null" json:"jumlah_pinjaman"`
	BungaOptionID       *uint        `gorm:"index" json:"bunga_option_id"`                   // References bunga_options table (nullable for backward compatibility)
	BungaPersen         float64      `gorm:"type:decimal(5,2);not null" json:"bunga_persen"` // Copied from selected option for historical record
	LamaBulan           int          `gorm:"not null" json:"lama_bulan"`
	JumlahAngsuran      float64      `gorm:"type:decimal(15,2);not null" json:"jumlah_angsuran"`
	SisaAngsuran        int          `gorm:"not null" json:"sisa_angsuran"`
	Status              string       `gorm:"type:varchar(20);check:status IN ('proses', 'disetujui', 'lunas', 'macet')" json:"status"`
	NoRekeningPencairan string       `gorm:"type:varchar(50)" json:"no_rekening_pencairan"` // Account number for loan disbursement
	BankName            string       `gorm:"type:varchar(100)" json:"bank_name"`            // Bank name for disbursement
	User                User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BungaOption         *BungaOption `gorm:"foreignKey:BungaOptionID" json:"bunga_option,omitempty"`
}

// TableName specifies the table name for Pinjaman model
func (Pinjaman) TableName() string {
	return "pinjaman"
}
