package repository

import (
	"koperasi-service/internal/model"
	"time"

	"gorm.io/gorm"
)

type AuditTrailRepository interface {
	Create(audit *model.AuditTrail) error
	GetByID(id uint) (*model.AuditTrail, error)
	GetAll(filters AuditTrailFilters) ([]model.AuditTrail, int64, error)
	GetByUser(userID uint, limit, offset int) ([]model.AuditTrail, error)
	GetByAction(action string, limit, offset int) ([]model.AuditTrail, error)
	GetByTable(table string, limit, offset int) ([]model.AuditTrail, error)
	GetByDateRange(startDate, endDate time.Time, limit, offset int) ([]model.AuditTrail, error)
	GetUserActivity(userID uint, startDate, endDate time.Time) ([]model.AuditTrail, error)
	GetSystemActivity(startDate, endDate time.Time) ([]model.AuditTrail, error)
}

type AuditTrailFilters struct {
	UserID      *uint      `json:"user_id"`
	Action      string     `json:"action"`
	EntityTable string     `json:"entity_table"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	IPAddress   string     `json:"ip_address"`
	Limit       int        `json:"limit"`
	Offset      int        `json:"offset"`
}

type auditTrailRepository struct {
	db *gorm.DB
}

func NewAuditTrailRepository(db *gorm.DB) AuditTrailRepository {
	return &auditTrailRepository{db: db}
}

func (r *auditTrailRepository) Create(audit *model.AuditTrail) error {
	return r.db.Create(audit).Error
}

func (r *auditTrailRepository) GetByID(id uint) (*model.AuditTrail, error) {
	var audit model.AuditTrail
	err := r.db.Preload("User").First(&audit, id).Error
	if err != nil {
		return nil, err
	}
	return &audit, nil
}

func (r *auditTrailRepository) GetAll(filters AuditTrailFilters) ([]model.AuditTrail, int64, error) {
	var audits []model.AuditTrail
	var total int64

	query := r.db.Model(&model.AuditTrail{})

	// Apply filters
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.Action != "" {
		query = query.Where("action = ?", filters.Action)
	}
	if filters.EntityTable != "" {
		query = query.Where("entity_table = ?", filters.EntityTable)
	}
	if filters.StartDate != nil {
		query = query.Where("timestamp >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("timestamp <= ?", *filters.EndDate)
	}
	if filters.IPAddress != "" {
		query = query.Where("ip_address = ?", filters.IPAddress)
	}

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	if filters.Limit <= 0 {
		filters.Limit = 50 // Default limit
	}

	err = query.Preload("User").
		Order("timestamp DESC").
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&audits).Error

	return audits, total, err
}

func (r *auditTrailRepository) GetByUser(userID uint, limit, offset int) ([]model.AuditTrail, error) {
	var audits []model.AuditTrail
	err := r.db.Where("user_id = ?", userID).
		Preload("User").
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&audits).Error
	return audits, err
}

func (r *auditTrailRepository) GetByAction(action string, limit, offset int) ([]model.AuditTrail, error) {
	var audits []model.AuditTrail
	err := r.db.Where("action = ?", action).
		Preload("User").
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&audits).Error
	return audits, err
}

func (r *auditTrailRepository) GetByTable(table string, limit, offset int) ([]model.AuditTrail, error) {
	var audits []model.AuditTrail
	err := r.db.Where("entity_table = ?", table).
		Preload("User").
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&audits).Error
	return audits, err
}

func (r *auditTrailRepository) GetByDateRange(startDate, endDate time.Time, limit, offset int) ([]model.AuditTrail, error) {
	var audits []model.AuditTrail
	err := r.db.Where("timestamp BETWEEN ? AND ?", startDate, endDate).
		Preload("User").
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&audits).Error
	return audits, err
}

func (r *auditTrailRepository) GetUserActivity(userID uint, startDate, endDate time.Time) ([]model.AuditTrail, error) {
	var audits []model.AuditTrail
	err := r.db.Where("user_id = ? AND timestamp BETWEEN ? AND ?", userID, startDate, endDate).
		Preload("User").
		Order("timestamp DESC").
		Find(&audits).Error
	return audits, err
}

func (r *auditTrailRepository) GetSystemActivity(startDate, endDate time.Time) ([]model.AuditTrail, error) {
	var audits []model.AuditTrail
	err := r.db.Where("timestamp BETWEEN ? AND ?", startDate, endDate).
		Preload("User").
		Order("timestamp DESC").
		Find(&audits).Error
	return audits, err
}

// TransactionHistoryRepository handles transaction history data operations
type TransactionHistoryRepository interface {
	Create(transaction *model.TransactionHistory) error
	GetByID(id uint) (*model.TransactionHistory, error)
	GetAll(filters TransactionHistoryFilters) ([]model.TransactionHistory, int64, error)
	GetByUser(userID uint, limit, offset int) ([]model.TransactionHistory, error)
	GetByType(transactionType string, limit, offset int) ([]model.TransactionHistory, error)
	GetByDateRange(startDate, endDate time.Time, limit, offset int) ([]model.TransactionHistory, error)
	GetUserTransactions(userID uint, startDate, endDate time.Time) ([]model.TransactionHistory, error)
	GetFinancialSummary(startDate, endDate time.Time) (map[string]interface{}, error)
	UpdateStatus(id uint, status string, verifiedBy uint) error
}

type TransactionHistoryFilters struct {
	UserID          *uint      `json:"user_id"`
	TransactionType string     `json:"transaction_type"`
	Status          string     `json:"status"`
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	MinAmount       *float64   `json:"min_amount"`
	MaxAmount       *float64   `json:"max_amount"`
	VerifiedBy      *uint      `json:"verified_by"`
	Limit           int        `json:"limit"`
	Offset          int        `json:"offset"`
}

type transactionHistoryRepository struct {
	db *gorm.DB
}

func NewTransactionHistoryRepository(db *gorm.DB) TransactionHistoryRepository {
	return &transactionHistoryRepository{db: db}
}

func (r *transactionHistoryRepository) Create(transaction *model.TransactionHistory) error {
	return r.db.Create(transaction).Error
}

func (r *transactionHistoryRepository) GetByID(id uint) (*model.TransactionHistory, error) {
	var transaction model.TransactionHistory
	err := r.db.Preload("User").Preload("VerifiedByUser").First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionHistoryRepository) GetAll(filters TransactionHistoryFilters) ([]model.TransactionHistory, int64, error) {
	var transactions []model.TransactionHistory
	var total int64

	query := r.db.Model(&model.TransactionHistory{})

	// Apply filters
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.TransactionType != "" {
		query = query.Where("transaction_type = ?", filters.TransactionType)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.StartDate != nil {
		query = query.Where("transaction_date >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("transaction_date <= ?", *filters.EndDate)
	}
	if filters.MinAmount != nil {
		query = query.Where("amount >= ?", *filters.MinAmount)
	}
	if filters.MaxAmount != nil {
		query = query.Where("amount <= ?", *filters.MaxAmount)
	}
	if filters.VerifiedBy != nil {
		query = query.Where("verified_by = ?", *filters.VerifiedBy)
	}

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	if filters.Limit <= 0 {
		filters.Limit = 50 // Default limit
	}

	err = query.Preload("User").Preload("VerifiedByUser").
		Order("transaction_date DESC").
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&transactions).Error

	return transactions, total, err
}

func (r *transactionHistoryRepository) GetByUser(userID uint, limit, offset int) ([]model.TransactionHistory, error) {
	var transactions []model.TransactionHistory
	err := r.db.Where("user_id = ?", userID).
		Preload("User").Preload("VerifiedByUser").
		Order("transaction_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionHistoryRepository) GetByType(transactionType string, limit, offset int) ([]model.TransactionHistory, error) {
	var transactions []model.TransactionHistory
	err := r.db.Where("transaction_type = ?", transactionType).
		Preload("User").Preload("VerifiedByUser").
		Order("transaction_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionHistoryRepository) GetByDateRange(startDate, endDate time.Time, limit, offset int) ([]model.TransactionHistory, error) {
	var transactions []model.TransactionHistory
	err := r.db.Where("transaction_date BETWEEN ? AND ?", startDate, endDate).
		Preload("User").Preload("VerifiedByUser").
		Order("transaction_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionHistoryRepository) GetUserTransactions(userID uint, startDate, endDate time.Time) ([]model.TransactionHistory, error) {
	var transactions []model.TransactionHistory
	err := r.db.Where("user_id = ? AND transaction_date BETWEEN ? AND ?", userID, startDate, endDate).
		Preload("User").Preload("VerifiedByUser").
		Order("transaction_date DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionHistoryRepository) GetFinancialSummary(startDate, endDate time.Time) (map[string]interface{}, error) {
	var summary struct {
		TotalSimpanan     float64 `gorm:"column:total_simpanan"`
		TotalPinjaman     float64 `gorm:"column:total_pinjaman"`
		TotalAngsuran     float64 `gorm:"column:total_angsuran"`
		TotalSHU          float64 `gorm:"column:total_shu"`
		TotalTransactions int64   `gorm:"column:total_transactions"`
	}

	err := r.db.Model(&model.TransactionHistory{}).
		Select(`
			COALESCE(SUM(CASE WHEN transaction_type = 'SIMPANAN' THEN amount END), 0) as total_simpanan,
			COALESCE(SUM(CASE WHEN transaction_type = 'PINJAMAN' THEN amount END), 0) as total_pinjaman,
			COALESCE(SUM(CASE WHEN transaction_type = 'ANGSURAN' THEN amount END), 0) as total_angsuran,
			COALESCE(SUM(CASE WHEN transaction_type = 'SHU' THEN amount END), 0) as total_shu,
			COUNT(*) as total_transactions
		`).
		Where("transaction_date BETWEEN ? AND ?", startDate, endDate).
		Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"total_simpanan":     summary.TotalSimpanan,
		"total_pinjaman":     summary.TotalPinjaman,
		"total_angsuran":     summary.TotalAngsuran,
		"total_shu":          summary.TotalSHU,
		"total_transactions": summary.TotalTransactions,
		"period_start":       startDate,
		"period_end":         endDate,
	}

	return result, nil
}

func (r *transactionHistoryRepository) UpdateStatus(id uint, status string, verifiedBy uint) error {
	now := time.Now()
	return r.db.Model(&model.TransactionHistory{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"verified_by": verifiedBy,
			"verified_at": &now,
		}).Error
}
