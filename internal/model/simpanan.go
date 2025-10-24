package model

import "gorm.io/gorm"

// Simpanan represents a savings wallet for each user with three types
type Simpanan struct {
	gorm.Model
	UserID      uint
	Type        string  // "pokok", "wajib", "sukarela"
	Balance     float64 // Current balance in the wallet
	Description string
}

// SimpananTransaction represents top-up or adjustment transactions
type SimpananTransaction struct {
	gorm.Model
	SimpananID   uint // Reference to the simpanan wallet
	Simpanan     Simpanan
	Type         string  // "topup", "adjustment"
	Amount       float64 // Amount of transaction (positive for topup, negative for deduction)
	Description  string
	Status       string // "pending", "verified", "rejected"
	VerifiedByID *uint  // Admin who verified the transaction
	VerifiedBy   *User  `gorm:"foreignKey:VerifiedByID"`
	VerifiedAt   *gorm.DeletedAt
}
