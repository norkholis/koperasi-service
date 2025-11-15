package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"koperasi-service/internal/model"
	"koperasi-service/internal/repository"
	"time"
)

type AuditTrailService interface {
	CreateAuditLog(userID uint, action, entityTable string, recordID uint, oldValues, newValues interface{}, ipAddress, userAgent, description string) error
	GetAuditTrails(userID uint, filters repository.AuditTrailFilters) ([]model.AuditTrail, int64, error)
	GetAuditTrailByID(userID uint, id uint) (*model.AuditTrail, error)
	GetUserActivity(userID uint, targetUserID uint, startDate, endDate time.Time) ([]model.AuditTrail, error)
	GetSystemActivity(userID uint, startDate, endDate time.Time) ([]model.AuditTrail, error)
	GetAuditSummary(userID uint, startDate, endDate time.Time) (map[string]interface{}, error)
}

type auditTrailService struct {
	auditRepo repository.AuditTrailRepository
	userRepo  *repository.UserRepository
}

func NewAuditTrailService(auditRepo repository.AuditTrailRepository, userRepo *repository.UserRepository) AuditTrailService {
	return &auditTrailService{
		auditRepo: auditRepo,
		userRepo:  userRepo,
	}
}

func (s *auditTrailService) CreateAuditLog(userID uint, action, entityTable string, recordID uint, oldValues, newValues interface{}, ipAddress, userAgent, description string) error {
	var oldValuesJSON, newValuesJSON string

	if oldValues != nil {
		oldBytes, err := json.Marshal(oldValues)
		if err != nil {
			return fmt.Errorf("failed to marshal old values: %v", err)
		}
		oldValuesJSON = string(oldBytes)
	}

	if newValues != nil {
		newBytes, err := json.Marshal(newValues)
		if err != nil {
			return fmt.Errorf("failed to marshal new values: %v", err)
		}
		newValuesJSON = string(newBytes)
	}

	audit := &model.AuditTrail{
		UserID:      userID,
		Action:      action,
		EntityTable: entityTable,
		RecordID:    recordID,
		OldValues:   oldValuesJSON,
		NewValues:   newValuesJSON,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Description: description,
		Timestamp:   time.Now(),
	}

	return s.auditRepo.Create(audit)
}

func (s *auditTrailService) GetAuditTrails(userID uint, filters repository.AuditTrailFilters) ([]model.AuditTrail, int64, error) {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return nil, 0, errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return nil, 0, errors.New("only admin and super admin can access audit trails")
	}

	return s.auditRepo.GetAll(filters)
}

func (s *auditTrailService) GetAuditTrailByID(userID uint, id uint) (*model.AuditTrail, error) {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return nil, errors.New("only admin and super admin can access audit trails")
	}

	return s.auditRepo.GetByID(id)
}

func (s *auditTrailService) GetUserActivity(userID uint, targetUserID uint, startDate, endDate time.Time) ([]model.AuditTrail, error) {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return nil, errors.New("only admin and super admin can access user activity")
	}

	return s.auditRepo.GetUserActivity(targetUserID, startDate, endDate)
}

func (s *auditTrailService) GetSystemActivity(userID uint, startDate, endDate time.Time) ([]model.AuditTrail, error) {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return nil, errors.New("only admin and super admin can access system activity")
	}

	return s.auditRepo.GetSystemActivity(startDate, endDate)
}

func (s *auditTrailService) GetAuditSummary(userID uint, startDate, endDate time.Time) (map[string]interface{}, error) {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return nil, errors.New("only admin and super admin can access audit summary")
	}

	// Get audit trails for the period
	filters := repository.AuditTrailFilters{
		StartDate: &startDate,
		EndDate:   &endDate,
		Limit:     1000, // Large limit for summary
	}

	audits, total, err := s.auditRepo.GetAll(filters)
	if err != nil {
		return nil, err
	}

	// Analyze audit data
	actionCounts := make(map[string]int)
	tableCounts := make(map[string]int)
	userCounts := make(map[uint]int)
	ipCounts := make(map[string]int)

	for _, audit := range audits {
		actionCounts[audit.Action]++
		tableCounts[audit.EntityTable]++
		userCounts[audit.UserID]++
		if audit.IPAddress != "" {
			ipCounts[audit.IPAddress]++
		}
	}

	summary := map[string]interface{}{
		"period_start":   startDate,
		"period_end":     endDate,
		"total_entries":  total,
		"action_summary": actionCounts,
		"table_summary":  tableCounts,
		"active_users":   len(userCounts),
		"unique_ips":     len(ipCounts),
		"top_actions":    getTopItems(actionCounts, 5),
		"top_tables":     getTopItems(tableCounts, 5),
		"top_ips":        getTopItems(ipCounts, 5),
	}

	return summary, nil
}

// TransactionHistoryService handles transaction history operations
type TransactionHistoryService interface {
	CreateTransactionRecord(userID uint, transactionType, referenceTable string, referenceID uint, amount, balanceBefore, balanceAfter float64, status, description, metadata string) error
	GetTransactionHistory(userID uint, filters repository.TransactionHistoryFilters) ([]model.TransactionHistory, int64, error)
	GetTransactionByID(userID uint, id uint) (*model.TransactionHistory, error)
	GetUserTransactions(userID uint, targetUserID uint, startDate, endDate time.Time) ([]model.TransactionHistory, error)
	GetFinancialSummary(userID uint, startDate, endDate time.Time) (map[string]interface{}, error)
	UpdateTransactionStatus(userID uint, id uint, status string) error
	GenerateFinancialReport(userID uint, reportType string, startDate, endDate time.Time) (map[string]interface{}, error)
}

type transactionHistoryService struct {
	transactionRepo repository.TransactionHistoryRepository
	userRepo        *repository.UserRepository
}

func NewTransactionHistoryService(transactionRepo repository.TransactionHistoryRepository, userRepo *repository.UserRepository) TransactionHistoryService {
	return &transactionHistoryService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

func (s *transactionHistoryService) CreateTransactionRecord(userID uint, transactionType, referenceTable string, referenceID uint, amount, balanceBefore, balanceAfter float64, status, description, metadata string) error {
	transaction := &model.TransactionHistory{
		UserID:          userID,
		TransactionType: transactionType,
		ReferenceTable:  referenceTable,
		ReferenceID:     referenceID,
		Amount:          amount,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		Status:          status,
		TransactionDate: time.Now(),
		Description:     description,
		Metadata:        metadata,
	}

	return s.transactionRepo.Create(transaction)
}

func (s *transactionHistoryService) GetTransactionHistory(userID uint, filters repository.TransactionHistoryFilters) ([]model.TransactionHistory, int64, error) {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return nil, 0, errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return nil, 0, errors.New("only admin and super admin can access transaction history")
	}

	return s.transactionRepo.GetAll(filters)
}

func (s *transactionHistoryService) GetTransactionByID(userID uint, id uint) (*model.TransactionHistory, error) {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return nil, errors.New("only admin and super admin can access transaction details")
	}

	return s.transactionRepo.GetByID(id)
}

func (s *transactionHistoryService) GetUserTransactions(userID uint, targetUserID uint, startDate, endDate time.Time) ([]model.TransactionHistory, error) {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return nil, errors.New("only admin and super admin can access user transactions")
	}

	return s.transactionRepo.GetUserTransactions(targetUserID, startDate, endDate)
}

func (s *transactionHistoryService) GetFinancialSummary(userID uint, startDate, endDate time.Time) (map[string]interface{}, error) {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return nil, errors.New("only admin and super admin can access financial summary")
	}

	return s.transactionRepo.GetFinancialSummary(startDate, endDate)
}

func (s *transactionHistoryService) UpdateTransactionStatus(userID uint, id uint, status string) error {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return errors.New("only admin and super admin can update transaction status")
	}

	return s.transactionRepo.UpdateStatus(id, status, userID)
}

func (s *transactionHistoryService) GenerateFinancialReport(userID uint, reportType string, startDate, endDate time.Time) (map[string]interface{}, error) {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return nil, errors.New("only admin and super admin can generate financial reports")
	}

	// Get financial summary
	summary, err := s.transactionRepo.GetFinancialSummary(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get transaction details by type
	filters := repository.TransactionHistoryFilters{
		StartDate: &startDate,
		EndDate:   &endDate,
		Limit:     1000,
	}

	transactions, total, err := s.transactionRepo.GetAll(filters)
	if err != nil {
		return nil, err
	}

	// Analyze transaction patterns
	typeBreakdown := make(map[string]float64)
	statusBreakdown := make(map[string]int)
	monthlyBreakdown := make(map[string]float64)

	for _, txn := range transactions {
		typeBreakdown[txn.TransactionType] += txn.Amount
		statusBreakdown[txn.Status]++

		month := txn.TransactionDate.Format("2006-01")
		monthlyBreakdown[month] += txn.Amount
	}

	report := map[string]interface{}{
		"report_type":       reportType,
		"period_start":      startDate,
		"period_end":        endDate,
		"generated_at":      time.Now(),
		"generated_by":      userID,
		"summary":           summary,
		"transaction_count": total,
		"type_breakdown":    typeBreakdown,
		"status_breakdown":  statusBreakdown,
		"monthly_breakdown": monthlyBreakdown,
	}

	return report, nil
}

// Helper function to get top items from a map
func getTopItems(counts map[string]int, limit int) []map[string]interface{} {
	type item struct {
		key   string
		count int
	}

	items := make([]item, 0, len(counts))
	for k, v := range counts {
		items = append(items, item{key: k, count: v})
	}

	// Simple sort by count (descending)
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].count > items[i].count {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	result := make([]map[string]interface{}, 0, limit)
	for i := 0; i < limit && i < len(items); i++ {
		result = append(result, map[string]interface{}{
			"item":  items[i].key,
			"count": items[i].count,
		})
	}

	return result
}
