package service

import (
	"errors"
	"koperasi-service/internal/model"
	"koperasi-service/internal/repository"
	"time"
)

// SHUService handles business logic for SHU calculations and management
type SHUService struct {
	repo *repository.SHUTahunanRepository
}

// NewSHUService creates a new service instance
func NewSHUService(repo *repository.SHUTahunanRepository) *SHUService {
	return &SHUService{repo: repo}
}

// SHU calculation constants (can be made configurable later)
const (
	DefaultPersenJasaModal = 25.0 // 25%
	DefaultPersenJasaUsaha = 30.0 // 30%
)

// GenerateReport calculates and generates SHU report for a specific year
func (s *SHUService) GenerateReport(requestorRole string, tahun int, totalSHUKoperasi float64) (*model.SHUReport, error) {
	// Only admin and super_admin can generate SHU reports
	if requestorRole != "admin" && requestorRole != "super_admin" {
		return nil, errors.New("forbidden")
	}

	// Get all required data for calculation
	totalSimpananAll, err := s.repo.GetTotalSimpananByYear(tahun)
	if err != nil {
		return nil, err
	}

	totalPenjualanAll, err := s.repo.GetTotalPenjualanByYear(tahun)
	if err != nil {
		return nil, err
	}

	userSimpanan, err := s.repo.GetSimpananByUserAndYear(tahun)
	if err != nil {
		return nil, err
	}

	userPenjualan, err := s.repo.GetPenjualanByUserAndYear(tahun)
	if err != nil {
		return nil, err
	}

	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	// Calculate SHU for each member
	var detailAnggota []model.SHUAnggota

	for _, user := range users {
		simpananAnggota := userSimpanan[user.ID]
		penjualanAnggota := userPenjualan[user.ID]

		// Calculate Jasa Modal Anggota (JMA)
		// JMA = (Simpanan anggota / Total simpanan koperasi) x % Jasa Modal x Total SHU Koperasi
		var jasaModal float64
		if totalSimpananAll > 0 {
			jasaModal = (simpananAnggota / totalSimpananAll) * (DefaultPersenJasaModal / 100) * totalSHUKoperasi
		}

		// Calculate Jasa Usaha Anggota (JUA)
		// JUA = (Total penjualan anggota / Total penjualan koperasi) x % Jasa Usaha x Total SHU Koperasi
		var jasaUsaha float64
		if totalPenjualanAll > 0 {
			jasaUsaha = (penjualanAnggota / totalPenjualanAll) * (DefaultPersenJasaUsaha / 100) * totalSHUKoperasi
		}

		// Total SHU Anggota = JMA + JUA
		totalSHUAnggota := jasaModal + jasaUsaha

		// Only include members who have some activity (simpanan or penjualan)
		if simpananAnggota > 0 || penjualanAnggota > 0 {
			detailAnggota = append(detailAnggota, model.SHUAnggota{
				UserID:          user.ID,
				Email:           user.Email,
				TotalSimpanan:   simpananAnggota,
				TotalPenjualan:  penjualanAnggota,
				JasaModal:       jasaModal,
				JasaUsaha:       jasaUsaha,
				TotalSHUAnggota: totalSHUAnggota,
			})
		}
	}

	report := &model.SHUReport{
		Tahun:              tahun,
		TotalSHUKoperasi:   totalSHUKoperasi,
		PersenJasaModal:    DefaultPersenJasaModal,
		PersenJasaUsaha:    DefaultPersenJasaUsaha,
		TotalSimpananAll:   totalSimpananAll,
		TotalPenjualanAll:  totalPenjualanAll,
		TanggalHitung:      time.Now(),
		DetailAnggota:      detailAnggota,
	}

	return report, nil
}

// SaveSHU saves the SHU calculation as a record
func (s *SHUService) SaveSHU(requestorRole string, tahun int, totalSHU float64, status string) (*model.SHUTahunan, error) {
	// Only admin and super_admin can save SHU
	if requestorRole != "admin" && requestorRole != "super_admin" {
		return nil, errors.New("forbidden")
	}

	// Check if SHU for this year already exists
	existing, _ := s.repo.GetByTahun(tahun)
	if existing != nil {
		return nil, errors.New("SHU for this year already exists")
	}

	// Validate status
	if status != "draft" && status != "final" {
		status = "draft"
	}

	shu := &model.SHUTahunan{
		Tahun:         tahun,
		TotalSHU:      totalSHU,
		TanggalHitung: time.Now(),
		Status:        status,
	}

	if err := s.repo.Create(shu); err != nil {
		return nil, err
	}

	return shu, nil
}

// List returns all SHU records
func (s *SHUService) List(requestorRole string) ([]model.SHUTahunan, error) {
	// Only admin and super_admin can view SHU records
	if requestorRole != "admin" && requestorRole != "super_admin" {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetAll()
}

// Get returns SHU record by id
func (s *SHUService) Get(requestorRole string, id uint) (*model.SHUTahunan, error) {
	// Only admin and super_admin can view SHU records
	if requestorRole != "admin" && requestorRole != "super_admin" {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetByID(id)
}

// Update modifies an existing SHU record
func (s *SHUService) Update(requestorRole string, id uint, payload *model.SHUTahunan) (*model.SHUTahunan, error) {
	// Only admin and super_admin can update SHU records
	if requestorRole != "admin" && requestorRole != "super_admin" {
		return nil, errors.New("forbidden")
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update allowed fields
	if payload.TotalSHU > 0 {
		existing.TotalSHU = payload.TotalSHU
	}
	if payload.Status != "" {
		// Validate status
		if payload.Status == "draft" || payload.Status == "final" {
			existing.Status = payload.Status
		}
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete removes a SHU record
func (s *SHUService) Delete(requestorRole string, id uint) error {
	// Only admin and super_admin can delete SHU records
	if requestorRole != "admin" && requestorRole != "super_admin" {
		return errors.New("forbidden")
	}

	// Check if record exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

// GetByTahun returns SHU record by year
func (s *SHUService) GetByTahun(requestorRole string, tahun int) (*model.SHUTahunan, error) {
	// Only admin and super_admin can view SHU records
	if requestorRole != "admin" && requestorRole != "super_admin" {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetByTahun(tahun)
}