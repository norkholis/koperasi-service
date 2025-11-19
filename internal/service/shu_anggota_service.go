package service

import (
	"errors"
	"fmt"
	"koperasi-service/internal/model"
	"koperasi-service/internal/repository"
)

// SHUAnggotaService handles business logic for SHU Anggota operations
type SHUAnggotaService struct {
	repo           *repository.SHUAnggotaRepository
	shuRepo        *repository.SHUTahunanRepository
	transactionSvc TransactionHistoryService
}

// NewSHUAnggotaService creates a new service instance
func NewSHUAnggotaService(repo *repository.SHUAnggotaRepository, shuRepo *repository.SHUTahunanRepository, transactionSvc TransactionHistoryService) *SHUAnggotaService {
	return &SHUAnggotaService{
		repo:           repo,
		shuRepo:        shuRepo,
		transactionSvc: transactionSvc,
	}
}

// SaveUserSHU saves the calculated SHU for a specific user
func (s *SHUAnggotaService) SaveUserSHU(requestorRole string, requestorUserID uint, targetUserID uint, tahun int) (*model.SHUAnggotaRecord, error) {
	// Admin and super_admin can save SHU for any user
	// Regular users can only save SHU for themselves
	if requestorRole != "admin" && requestorRole != "super_admin" {
		if requestorUserID != targetUserID {
			return nil, errors.New("forbidden")
		}
	}

	// Get the SHU record for the year
	shuRecord, err := s.shuRepo.GetByTahun(tahun)
	if err != nil {
		return nil, errors.New("SHU record not found for the specified year")
	}

	// Check if user SHU already exists
	exists, err := s.repo.CheckExists(shuRecord.ID, targetUserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("SHU record already exists for this user and year")
	}

	// Generate the SHU calculation for the user
	shuService := NewSHUService(s.shuRepo)
	shuAnggota, err := shuService.GenerateUserSHU(requestorRole, requestorUserID, targetUserID, tahun)
	if err != nil {
		return nil, err
	}

	// Create the SHU Anggota record
	shuAnggotaRecord := &model.SHUAnggotaRecord{
		SHUID:       shuRecord.ID,
		UserID:      targetUserID,
		JumlahModal: shuAnggota.JasaModal,
		JumlahUsaha: shuAnggota.JasaUsaha,
		SHUDiterima: shuAnggota.TotalSHUAnggota,
	}

	if err := s.repo.Create(shuAnggotaRecord); err != nil {
		return nil, err
	}

	// Record transaction history for SHU distribution
	if s.transactionSvc != nil {
		metadata := fmt.Sprintf(`{"tahun":%d,"jasa_modal":%.2f,"jasa_usaha":%.2f,"total_shu":%.2f}`, tahun, shuAnggota.JasaModal, shuAnggota.JasaUsaha, shuAnggota.TotalSHUAnggota)
		s.transactionSvc.CreateTransactionRecord(
			targetUserID,
			"SHU",
			"shu_anggota_records",
			shuAnggotaRecord.ID,
			shuAnggota.TotalSHUAnggota,
			0,
			shuAnggota.TotalSHUAnggota,
			"COMPLETED",
			fmt.Sprintf("SHU distribution for year %d: %.2f", tahun, shuAnggota.TotalSHUAnggota),
			metadata,
		)
	}

	return shuAnggotaRecord, nil
}

// GetUserSHU retrieves saved SHU data for a user
func (s *SHUAnggotaService) GetUserSHU(requestorRole string, requestorUserID uint, targetUserID uint, tahun int) (*model.SHUAnggotaRecord, error) {
	// Admin and super_admin can view SHU for any user
	// Regular users can only view SHU for themselves
	if requestorRole != "admin" && requestorRole != "super_admin" {
		if requestorUserID != targetUserID {
			return nil, errors.New("forbidden")
		}
	}

	// Get the SHU record for the year
	shuRecord, err := s.shuRepo.GetByTahun(tahun)
	if err != nil {
		return nil, errors.New("SHU record not found for the specified year")
	}

	// Get the user's SHU record
	shuAnggotaRecord, err := s.repo.GetBySHUIDAndUserID(shuRecord.ID, targetUserID)
	if err != nil {
		return nil, errors.New("user SHU record not found")
	}

	return shuAnggotaRecord, nil
}

// GetUserSHUHistory retrieves all SHU history for a user
func (s *SHUAnggotaService) GetUserSHUHistory(requestorRole string, requestorUserID uint, targetUserID uint) ([]model.SHUAnggotaRecord, error) {
	// Admin and super_admin can view SHU history for any user
	// Regular users can only view SHU history for themselves
	if requestorRole != "admin" && requestorRole != "super_admin" {
		if requestorUserID != targetUserID {
			return nil, errors.New("forbidden")
		}
	}

	return s.repo.GetByUserID(targetUserID)
}

// List retrieves all SHU Anggota records (admin only)
func (s *SHUAnggotaService) List(requestorRole string) ([]model.SHUAnggotaRecord, error) {
	if requestorRole != "admin" && requestorRole != "super_admin" {
		return nil, errors.New("forbidden")
	}

	return s.repo.List(0, 100) // Default pagination
}

// GetBySHUID retrieves all user SHU records for a specific SHU (admin only)
func (s *SHUAnggotaService) GetBySHUID(requestorRole string, shuID uint) ([]model.SHUAnggotaRecord, error) {
	if requestorRole != "admin" && requestorRole != "super_admin" {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetBySHUID(shuID)
}

// Delete removes a SHU Anggota record (admin only)
func (s *SHUAnggotaService) Delete(requestorRole string, id uint) error {
	if requestorRole != "admin" && requestorRole != "super_admin" {
		return errors.New("forbidden")
	}

	// Check if record exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("SHU Anggota record not found")
	}

	return s.repo.Delete(id)
}
