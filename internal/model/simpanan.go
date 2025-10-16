package model

import "gorm.io/gorm"

type Simpanan struct {
	gorm.Model
	UserID      uint
	Type        string // "wajib" atau "sukarela"
	Amount      float64
	Description string
}
