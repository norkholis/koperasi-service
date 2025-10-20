package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email       string `gorm:"unique"`
	Password    string
	Name        string
	Address     string
	PhoneNumber string
	NIK         string `gorm:"unique"` // National Identity Number (unique)
	RoleID      uint
	Role        Role
	AdminID     *uint `gorm:"index"`              // References the admin who registered this user
	Admin       *User `gorm:"foreignKey:AdminID"` // The admin who registered this user
}

type Role struct {
	gorm.Model
	Name string `gorm:"unique"`
}
