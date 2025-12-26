package model

import "gorm.io/gorm"

// Simpanan represents a savings wallet for each user with three types
type Simpanan struct {
	gorm.Model
	UserID      uint    `gorm:"uniqueIndex:idx_user_wallet_type"`
	Type        string  `gorm:"uniqueIndex:idx_user_wallet_type"` // "pokok", "wajib", "sukarela"
	Balance     float64 // Current balance in the wallet
	Description string
	User        User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// SimpananTransaction represents top-up or adjustment transactions
type SimpananTransaction struct {
	gorm.Model
	SimpananID         uint // Reference to the simpanan wallet
	Simpanan           Simpanan
	Type               string  // "topup", "adjustment"
	Amount             float64 // Amount of transaction (positive for topup, negative for deduction)
	Description        string
	Status             string       // "pending", "verified", "rejected"
	ImageBuktiTransfer string       `gorm:"type:varchar(255)" json:"image_bukti_transfer"` // URL/path to transfer proof image
	BankAccountID      *uint        `json:"bank_account_id"`                               // Which bank account was used for topup
	BankAccount        *BankAccount `gorm:"foreignKey:BankAccountID" json:"bank_account,omitempty"`
	VerifiedByID       *uint        // Admin who verified the transaction
	VerifiedBy         *User        `gorm:"foreignKey:VerifiedByID"`
	VerifiedAt         *gorm.DeletedAt
}
