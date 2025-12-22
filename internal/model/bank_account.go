package model

import "gorm.io/gorm"

// BankAccount represents bank accounts for receiving payments/topups
type BankAccount struct {
	gorm.Model
	BankName      string `gorm:"type:varchar(100);not null" json:"bank_name"`
	AccountNumber string `gorm:"type:varchar(50);not null" json:"account_number"`
	AccountName   string `gorm:"type:varchar(100);not null" json:"account_name"`
	IsActive      bool   `gorm:"default:true" json:"is_active"`
	Description   string `gorm:"type:text" json:"description"` // e.g., "For topup", "For loan payments"
}

// TableName specifies the table name for BankAccount model
func (BankAccount) TableName() string {
	return "bank_accounts"
}
