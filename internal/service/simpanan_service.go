package service

import (
	"errors"
	"koperasi-service/internal/model"
	"koperasi-service/internal/repository"
	"time"

	"gorm.io/gorm"
)

// SimpananService contains business logic for Simpanan wallets.
type SimpananService struct {
	repo *repository.SimpananRepository
}

// NewSimpananService creates a new service instance.
func NewSimpananService(repo *repository.SimpananRepository) *SimpananService {
	return &SimpananService{repo: repo}
}

// InitializeUserWallets creates the three wallet types for a new user
func (s *SimpananService) InitializeUserWallets(userID uint) error {
	return s.repo.InitializeUserWallets(userID)
}

// GetUserWallets returns all wallet types for a user
func (s *SimpananService) GetUserWallets(userID uint, requestorID uint, requestorRole string) ([]model.Simpanan, error) {
	// Check permissions
	if requestorRole == "super_admin" || requestorRole == "admin" {
		// Admin can see any user's wallets
		return s.repo.GetUserWallets(userID)
	}

	// Members can only see their own wallets
	if requestorID != userID {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetUserWallets(userID)
}

// GetAllWallets returns all wallets (admin only)
func (s *SimpananService) GetAllWallets(requestorRole string) ([]model.Simpanan, error) {
	if requestorRole != "super_admin" && requestorRole != "admin" {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetAllWallets(0)
}

// TopupWallet creates a pending top-up transaction
func (s *SimpananService) TopupWallet(userID uint, walletType string, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// Valid wallet types
	validTypes := map[string]bool{"pokok": true, "wajib": true, "sukarela": true}
	if !validTypes[walletType] {
		return errors.New("invalid wallet type")
	}

	// Get or create wallet
	wallet, err := s.repo.GetWalletByUserAndType(userID, walletType)
	if err != nil {
		return errors.New("wallet not found")
	}

	// Create pending transaction
	transaction := &model.SimpananTransaction{
		SimpananID:  wallet.ID,
		Type:        "topup",
		Amount:      amount,
		Description: description,
		Status:      "pending",
	}

	return s.repo.CreateTransaction(transaction)
}

// VerifyTransaction verifies and processes a pending transaction (admin only)
func (s *SimpananService) VerifyTransaction(transactionID uint, adminID uint, adminRole string, approve bool) error {
	if adminRole != "super_admin" && adminRole != "admin" {
		return errors.New("forbidden")
	}

	// Get transaction
	transaction, err := s.repo.GetTransactionByID(transactionID)
	if err != nil {
		return err
	}

	if transaction.Status != "pending" {
		return errors.New("transaction already processed")
	}

	// Update transaction status
	if approve {
		transaction.Status = "verified"

		// Update wallet balance
		wallet, err := s.repo.GetWalletByID(transaction.SimpananID)
		if err != nil {
			return err
		}

		wallet.Balance += transaction.Amount
		if err := s.repo.UpdateWallet(wallet); err != nil {
			return err
		}
	} else {
		transaction.Status = "rejected"
	}

	transaction.VerifiedByID = &adminID
	now := gorm.DeletedAt{Time: time.Now(), Valid: true}
	transaction.VerifiedAt = &now

	return s.repo.UpdateTransaction(transaction)
}

// AdjustWalletBalance allows admin to directly adjust wallet balance
func (s *SimpananService) AdjustWalletBalance(walletID uint, amount float64, description string, adminID uint, adminRole string) error {
	if adminRole != "super_admin" && adminRole != "admin" {
		return errors.New("forbidden")
	}

	wallet, err := s.repo.GetWalletByID(walletID)
	if err != nil {
		return err
	}

	// Create adjustment transaction
	transaction := &model.SimpananTransaction{
		SimpananID:   wallet.ID,
		Type:         "adjustment",
		Amount:       amount,
		Description:  description,
		Status:       "verified",
		VerifiedByID: &adminID,
		VerifiedAt:   &gorm.DeletedAt{Time: time.Now(), Valid: true},
	}

	if err := s.repo.CreateTransaction(transaction); err != nil {
		return err
	}

	// Update wallet balance immediately
	wallet.Balance += amount
	if wallet.Balance < 0 {
		return errors.New("insufficient balance")
	}

	return s.repo.UpdateWallet(wallet)
}

// GetWalletTransactions returns transaction history for a wallet
func (s *SimpananService) GetWalletTransactions(walletID uint, requestorID uint, requestorRole string) ([]model.SimpananTransaction, error) {
	// Get wallet to check ownership
	wallet, err := s.repo.GetWalletByID(walletID)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if requestorRole == "super_admin" || requestorRole == "admin" {
		// Admin can see any wallet's transactions
	} else if wallet.UserID != requestorID {
		// Members can only see their own wallet transactions
		return nil, errors.New("forbidden")
	}

	return s.repo.GetTransactionsByWallet(walletID)
}

// GetPendingTransactions returns all pending transactions (admin only)
func (s *SimpananService) GetPendingTransactions(adminRole string) ([]model.SimpananTransaction, error) {
	if adminRole != "super_admin" && adminRole != "admin" {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetPendingTransactions()
}

// GetWalletDetail returns detailed wallet information
func (s *SimpananService) GetWalletDetail(walletID uint, requestorID uint, requestorRole string) (*model.Simpanan, error) {
	wallet, err := s.repo.GetWalletByID(walletID)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if requestorRole == "super_admin" || requestorRole == "admin" {
		return wallet, nil
	}

	// Members can only see their own wallets
	if wallet.UserID != requestorID {
		return nil, errors.New("forbidden")
	}

	return wallet, nil
}
