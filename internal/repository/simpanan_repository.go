package repository

import (
	"koperasi-service/internal/model"

	"gorm.io/gorm"
)

// SimpananRepository handles persistence for Simpanan wallets and transactions.
type SimpananRepository struct {
	db *gorm.DB
}

// NewSimpananRepository constructs a new repository instance.
func NewSimpananRepository(db *gorm.DB) *SimpananRepository {
	return &SimpananRepository{db: db}
}

// InitializeUserWallets creates the three wallet types for a new user
func (r *SimpananRepository) InitializeUserWallets(userID uint) error {
	walletTypes := []string{"pokok", "wajib", "sukarela"}

	for _, walletType := range walletTypes {
		wallet := &model.Simpanan{
			UserID:      userID,
			Type:        walletType,
			Balance:     0,
			Description: "Wallet " + walletType,
		}
		if err := r.db.Create(wallet).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetUserWallets returns all wallet types for a specific user
func (r *SimpananRepository) GetUserWallets(userID uint) ([]model.Simpanan, error) {
	var wallets []model.Simpanan
	if err := r.db.Where("user_id = ?", userID).Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

// GetWalletByUserAndType returns a specific wallet type for a user
func (r *SimpananRepository) GetWalletByUserAndType(userID uint, walletType string) (*model.Simpanan, error) {
	var wallet model.Simpanan
	if err := r.db.Where("user_id = ? AND type = ?", userID, walletType).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

// GetAllWallets returns all wallets; if userID > 0 it filters by user.
func (r *SimpananRepository) GetAllWallets(userID uint) ([]model.Simpanan, error) {
	var list []model.Simpanan
	q := r.db
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetWalletByID returns single wallet by id.
func (r *SimpananRepository) GetWalletByID(id uint) (*model.Simpanan, error) {
	var s model.Simpanan
	if err := r.db.First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// UpdateWallet persists changes to an existing wallet.
func (r *SimpananRepository) UpdateWallet(s *model.Simpanan) error {
	return r.db.Save(s).Error
}

// CreateTransaction creates a new simpanan transaction
func (r *SimpananRepository) CreateTransaction(tx *model.SimpananTransaction) error {
	return r.db.Create(tx).Error
}

// GetTransactionsByWallet returns all transactions for a specific wallet
func (r *SimpananRepository) GetTransactionsByWallet(simpananID uint) ([]model.SimpananTransaction, error) {
	var transactions []model.SimpananTransaction
	if err := r.db.Where("simpanan_id = ?", simpananID).Preload("VerifiedBy").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// GetTransactionByID returns a transaction by ID
func (r *SimpananRepository) GetTransactionByID(id uint) (*model.SimpananTransaction, error) {
	var tx model.SimpananTransaction
	if err := r.db.Preload("Simpanan").Preload("VerifiedBy").First(&tx, id).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}

// UpdateTransaction updates a transaction
func (r *SimpananRepository) UpdateTransaction(tx *model.SimpananTransaction) error {
	return r.db.Save(tx).Error
}

// GetPendingTransactions returns all pending transactions (for admin verification)
func (r *SimpananRepository) GetPendingTransactions() ([]model.SimpananTransaction, error) {
	var transactions []model.SimpananTransaction
	if err := r.db.Where("status = ?", "pending").Preload("Simpanan").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
