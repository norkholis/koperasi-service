package repository

import (
	"koperasi-service/internal/model"

	"gorm.io/gorm"
)

type BankAccountRepository struct {
	db *gorm.DB
}

func NewBankAccountRepository(db *gorm.DB) *BankAccountRepository {
	return &BankAccountRepository{db: db}
}

func (r *BankAccountRepository) Create(account *model.BankAccount) error {
	return r.db.Create(account).Error
}

func (r *BankAccountRepository) Update(account *model.BankAccount) error {
	return r.db.Save(account).Error
}

func (r *BankAccountRepository) Delete(id uint) error {
	return r.db.Delete(&model.BankAccount{}, id).Error
}

func (r *BankAccountRepository) FindByID(id uint) (*model.BankAccount, error) {
	var account model.BankAccount
	if err := r.db.First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *BankAccountRepository) GetAll() ([]model.BankAccount, error) {
	var accounts []model.BankAccount
	if err := r.db.Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *BankAccountRepository) GetActive() ([]model.BankAccount, error) {
	var accounts []model.BankAccount
	if err := r.db.Where("is_active = ?", true).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
