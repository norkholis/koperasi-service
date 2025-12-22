package service

import (
	"errors"
	"koperasi-service/internal/model"
	"koperasi-service/internal/repository"
)

type BankAccountService struct {
	repo *repository.BankAccountRepository
}

func NewBankAccountService(repo *repository.BankAccountRepository) *BankAccountService {
	return &BankAccountService{repo: repo}
}

func (s *BankAccountService) Create(account *model.BankAccount) error {
	if account.BankName == "" || account.AccountNumber == "" || account.AccountName == "" {
		return errors.New("bank name, account number, and account name are required")
	}
	return s.repo.Create(account)
}

func (s *BankAccountService) Update(id uint, account *model.BankAccount) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("bank account not found")
	}

	// Update fields
	if account.BankName != "" {
		existing.BankName = account.BankName
	}
	if account.AccountNumber != "" {
		existing.AccountNumber = account.AccountNumber
	}
	if account.AccountName != "" {
		existing.AccountName = account.AccountName
	}
	existing.Description = account.Description
	existing.IsActive = account.IsActive

	return s.repo.Update(existing)
}

func (s *BankAccountService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("bank account not found")
	}
	return s.repo.Delete(id)
}

func (s *BankAccountService) GetByID(id uint) (*model.BankAccount, error) {
	return s.repo.FindByID(id)
}

func (s *BankAccountService) GetAll() ([]model.BankAccount, error) {
	return s.repo.GetAll()
}

func (s *BankAccountService) GetActive() ([]model.BankAccount, error) {
	return s.repo.GetActive()
}
