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
	DefaultPersenSHUAnggota = 50.0 // 50% of total SHU goes to members
	DefaultPersenJasaModal  = 30.0 // 30% of member SHU for Jasa Modal
	DefaultPersenJasaUsaha  = 70.0 // 70% of member SHU for Jasa Usaha
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

	// Calculate SHU allocation for members (50% of total SHU)
	shuUntukAnggota := totalSHUKoperasi * (DefaultPersenSHUAnggota / 100)
	alokasiJasaModal := shuUntukAnggota * (DefaultPersenJasaModal / 100)
	alokasiJasaUsaha := shuUntukAnggota * (DefaultPersenJasaUsaha / 100)

	for _, user := range users {
		simpananAnggota := userSimpanan[user.ID]
		pinjamanAnggota := userPenjualan[user.ID] // This is actually loan data from Pinjaman table

		// Calculate Jasa Modal Anggota (JMA)
		// JMA = (Simpanan anggota / Total simpanan koperasi) × Alokasi Jasa Modal
		var jasaModal float64
		if totalSimpananAll > 0 {
			jasaModal = (simpananAnggota / totalSimpananAll) * alokasiJasaModal
		}

		// Calculate Jasa Usaha Anggota (JUA)
		// JUA = (Pinjaman anggota / Total pinjaman koperasi) × Alokasi Jasa Usaha
		var jasaUsaha float64
		if totalPenjualanAll > 0 {
			jasaUsaha = (pinjamanAnggota / totalPenjualanAll) * alokasiJasaUsaha
		}

		// Total SHU Anggota = JMA + JUA
		totalSHUAnggota := jasaModal + jasaUsaha

		// Only include members who have some activity (simpanan or pinjaman)
		if simpananAnggota > 0 || pinjamanAnggota > 0 {
			detailAnggota = append(detailAnggota, model.SHUAnggota{
				UserID:          user.ID,
				Email:           user.Email,
				TotalSimpanan:   simpananAnggota,
				TotalPenjualan:  pinjamanAnggota, // Keep field name for compatibility but this is loan amount
				JasaModal:       jasaModal,
				JasaUsaha:       jasaUsaha,
				TotalSHUAnggota: totalSHUAnggota,
			})
		}
	}

	report := &model.SHUReport{
		Tahun:             tahun,
		TotalSHUKoperasi:  totalSHUKoperasi,
		PersenJasaModal:   DefaultPersenJasaModal,
		PersenJasaUsaha:   DefaultPersenJasaUsaha,
		TotalSimpananAll:  totalSimpananAll,
		TotalPenjualanAll: totalPenjualanAll, // Keep field name for compatibility but this is total loans
		TanggalHitung:     time.Now(),
		DetailAnggota:     detailAnggota,
	}

	return report, nil
}

// GenerateReportWithExpenses calculates SHU automatically based on income and expenses
func (s *SHUService) GenerateReportWithExpenses(requestorRole string, tahun int, bebanOperasional, bebanNonOperasional, bebanPajak float64) (*model.SHUReport, error) {
	// Only admin and super_admin can generate SHU reports
	if requestorRole != "admin" && requestorRole != "super_admin" {
		return nil, errors.New("forbidden")
	}

	// Calculate pendapatan (income) automatically
	pendapatanOperasional, err := s.repo.GetPendapatanOperasionalByYear(tahun)
	if err != nil {
		return nil, err
	}

	pendapatanNonOperasional, err := s.repo.GetPendapatanNonOperasionalByYear(tahun)
	if err != nil {
		return nil, err
	}

	// Calculate total SHU koperasi using the formula:
	// SHU Total = (Pendapatan Operasional + Pendapatan Non-Operasional) - (Beban Operasional + Beban Non-Operasional + Beban Pajak)
	totalPendapatan := pendapatanOperasional + pendapatanNonOperasional
	totalBeban := bebanOperasional + bebanNonOperasional + bebanPajak
	totalSHUKoperasi := totalPendapatan - totalBeban

	// Ensure SHU is not negative
	if totalSHUKoperasi < 0 {
		totalSHUKoperasi = 0
	}

	// Get all required data for member calculations
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

	// Calculate SHU allocation for members (50% of total SHU)
	shuUntukAnggota := totalSHUKoperasi * (DefaultPersenSHUAnggota / 100)
	alokasiJasaModal := shuUntukAnggota * (DefaultPersenJasaModal / 100)
	alokasiJasaUsaha := shuUntukAnggota * (DefaultPersenJasaUsaha / 100)

	for _, user := range users {
		simpananAnggota := userSimpanan[user.ID]
		pinjamanAnggota := userPenjualan[user.ID] // This is actually loan data from Pinjaman table

		// Calculate Jasa Modal Anggota (JMA)
		// JMA = (Simpanan anggota / Total simpanan koperasi) × Alokasi Jasa Modal
		var jasaModal float64
		if totalSimpananAll > 0 {
			jasaModal = (simpananAnggota / totalSimpananAll) * alokasiJasaModal
		}

		// Calculate Jasa Usaha Anggota (JUA)
		// JUA = (Pinjaman anggota / Total pinjaman koperasi) × Alokasi Jasa Usaha
		var jasaUsaha float64
		if totalPenjualanAll > 0 {
			jasaUsaha = (pinjamanAnggota / totalPenjualanAll) * alokasiJasaUsaha
		}

		// Total SHU Anggota = JMA + JUA
		totalSHUAnggota := jasaModal + jasaUsaha

		// Only include members who have some activity (simpanan or pinjaman)
		if simpananAnggota > 0 || pinjamanAnggota > 0 {
			detailAnggota = append(detailAnggota, model.SHUAnggota{
				UserID:          user.ID,
				Email:           user.Email,
				TotalSimpanan:   simpananAnggota,
				TotalPenjualan:  pinjamanAnggota, // Keep field name for compatibility but this is loan amount
				JasaModal:       jasaModal,
				JasaUsaha:       jasaUsaha,
				TotalSHUAnggota: totalSHUAnggota,
			})
		}
	}

	report := &model.SHUReport{
		Tahun:                    tahun,
		PendapatanOperasional:    pendapatanOperasional,
		PendapatanNonOperasional: pendapatanNonOperasional,
		BebanOperasional:         bebanOperasional,
		BebanNonOperasional:      bebanNonOperasional,
		BebanPajak:               bebanPajak,
		TotalSHUKoperasi:         totalSHUKoperasi,
		PersenJasaModal:          DefaultPersenJasaModal,
		PersenJasaUsaha:          DefaultPersenJasaUsaha,
		TotalSimpananAll:         totalSimpananAll,
		TotalPenjualanAll:        totalPenjualanAll,
		TanggalHitung:            time.Now(),
		DetailAnggota:            detailAnggota,
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

// SaveSHUWithExpenses saves the automated SHU calculation with detailed income and expense information
func (s *SHUService) SaveSHUWithExpenses(requestorRole string, tahun int, pendapatanOperasional, pendapatanNonOperasional, bebanOperasional, bebanNonOperasional, bebanPajak, totalSHU float64, status string) (*model.SHUTahunan, error) {
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
		Tahun:                    tahun,
		PendapatanOperasional:    pendapatanOperasional,
		PendapatanNonOperasional: pendapatanNonOperasional,
		BebanOperasional:         bebanOperasional,
		BebanNonOperasional:      bebanNonOperasional,
		BebanPajak:               bebanPajak,
		TotalSHU:                 totalSHU,
		TanggalHitung:            time.Now(),
		Status:                   status,
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

// GenerateUserSHU generates SHU report for a specific user
func (s *SHUService) GenerateUserSHU(requestorRole string, requestorUserID uint, targetUserID uint, tahun int) (*model.SHUAnggota, error) {
	// Admin and super_admin can generate SHU for any user
	// Regular users can only generate SHU for themselves
	if requestorRole != "admin" && requestorRole != "super_admin" {
		if requestorUserID != targetUserID {
			return nil, errors.New("forbidden")
		}
	}

	// Get existing SHU record for the year to get total SHU
	shuRecord, err := s.repo.GetByTahun(tahun)
	if err != nil {
		return nil, errors.New("SHU record not found for the specified year. Please generate SHU report first")
	}

	totalSHUKoperasi := shuRecord.TotalSHU

	// Get total simpanan and penjualan for the cooperative
	totalSimpananAll, err := s.repo.GetTotalSimpananByYear(tahun)
	if err != nil {
		return nil, err
	}

	totalPenjualanAll, err := s.repo.GetTotalPenjualanByYear(tahun)
	if err != nil {
		return nil, err
	}

	// Get user-specific data
	userSimpanan, err := s.repo.GetSimpananByUserAndYear(tahun)
	if err != nil {
		return nil, err
	}

	userPenjualan, err := s.repo.GetPenjualanByUserAndYear(tahun)
	if err != nil {
		return nil, err
	}

	// Get all users to find user email
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	// Find the specific user's email
	var userEmail string
	for _, user := range users {
		if user.ID == targetUserID {
			userEmail = user.Email
			break
		}
	}

	// If user not found, return error
	if userEmail == "" {
		return nil, errors.New("user not found")
	}

	// Get user's simpanan and pinjaman data
	targetUserSimpanan := userSimpanan[targetUserID]  // Will be 0 if not found
	targetUserPinjaman := userPenjualan[targetUserID] // Will be 0 if not found (this is actually loan data)

	// Calculate SHU allocation for members (50% of total SHU)
	shuUntukAnggota := totalSHUKoperasi * (DefaultPersenSHUAnggota / 100)
	alokasiJasaModal := shuUntukAnggota * (DefaultPersenJasaModal / 100)
	alokasiJasaUsaha := shuUntukAnggota * (DefaultPersenJasaUsaha / 100)

	// Calculate SHU for this specific user
	var jasaModal, jasaUsaha float64

	if totalSimpananAll > 0 {
		jasaModal = (targetUserSimpanan / totalSimpananAll) * alokasiJasaModal
	}

	if totalPenjualanAll > 0 {
		jasaUsaha = (targetUserPinjaman / totalPenjualanAll) * alokasiJasaUsaha
	}

	totalSHUAnggota := jasaModal + jasaUsaha

	shuAnggota := &model.SHUAnggota{
		UserID:          targetUserID,
		Email:           userEmail,
		TotalSimpanan:   targetUserSimpanan,
		TotalPenjualan:  targetUserPinjaman, // Keep field name for compatibility but this is loan amount
		JasaModal:       jasaModal,
		JasaUsaha:       jasaUsaha,
		TotalSHUAnggota: totalSHUAnggota,
	}

	return shuAnggota, nil
}
