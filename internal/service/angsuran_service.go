package service

import (
	"errors"
	"koperasi-service/internal/model"
	"koperasi-service/internal/repository"
	"time"
)

// AngsuranService handles business logic for Angsuran with role constraints
type AngsuranService struct {
	repo         *repository.AngsuranRepository
	pinjamanRepo *repository.PinjamanRepository
	userRepo     *repository.UserRepository
}

// NewAngsuranService creates a new service instance
func NewAngsuranService(repo *repository.AngsuranRepository, pinjamanRepo *repository.PinjamanRepository, userRepo *repository.UserRepository) *AngsuranService {
	return &AngsuranService{
		repo:         repo,
		pinjamanRepo: pinjamanRepo,
		userRepo:     userRepo,
	}
}

// Create adds a new Angsuran payment record
func (s *AngsuranService) Create(requestorID uint, requestorRole string, a *model.Angsuran) error {
	// Verify the pinjaman exists and access is allowed
	pinjaman, err := s.pinjamanRepo.GetByID(a.PinjamanID)
	if err != nil {
		return errors.New("pinjaman not found")
	}

	// Check access: super_admin can create for any, others only for their own loans
	if requestorRole != "super_admin" && pinjaman.UserID != requestorID {
		return errors.New("forbidden")
	}

	// Set default values
	if a.TanggalBayar.IsZero() {
		a.TanggalBayar = time.Now()
	}
	if a.Status == "" {
		a.Status = "proses"
	}
	if a.UserID == 0 {
		a.UserID = pinjaman.UserID
	}

	// Auto-set angsuran_ke if not provided
	if a.AngsuranKe == 0 {
		nextKe, err := s.repo.GetNextAngsuranKe(a.PinjamanID)
		if err != nil {
			return err
		}
		a.AngsuranKe = nextKe
	}

	// Calculate total if not provided
	if a.TotalBayar == 0 {
		a.TotalBayar = a.Pokok + a.Bunga + a.Denda
	}

	return s.repo.Create(a)
}

// List returns angsuran list filtered by access rules
func (s *AngsuranService) List(requestorID uint, requestorRole string, pinjamanID uint) ([]model.Angsuran, error) {
	if requestorRole == "super_admin" {
		// Super admin can see all angsuran
		return s.repo.GetAll(0, pinjamanID)
	} else if requestorRole == "admin" {
		// Admin can see angsuran for users they registered
		return s.repo.GetByAdminUserID(requestorID, pinjamanID)
	} else {
		// Members can only see their own angsuran
		return s.repo.GetAll(requestorID, pinjamanID)
	}
}

// Get returns Angsuran by id with access control
func (s *AngsuranService) Get(requestorID uint, requestorRole string, id uint) (*model.Angsuran, error) {
	a, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check access permissions based on role
	if requestorRole == "super_admin" {
		// Super admin can view any angsuran
		return a, nil
	} else if requestorRole == "admin" {
		// Admin can view angsuran for users they registered
		user, err := s.userRepo.FindByID(a.UserID)
		if err != nil {
			return nil, err
		}
		if user.AdminID == nil || *user.AdminID != requestorID {
			return nil, errors.New("forbidden")
		}
	} else {
		// Members can only view their own angsuran
		if a.UserID != requestorID {
			return nil, errors.New("forbidden")
		}
	}

	return a, nil
}

// Update modifies an existing Angsuran
func (s *AngsuranService) Update(requestorID uint, requestorRole string, id uint, payload *model.Angsuran) (*model.Angsuran, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Super admin can update any angsuran
	// Admin and member can only update their own angsuran
	if requestorRole != "super_admin" && existing.UserID != requestorID {
		return nil, errors.New("forbidden")
	}

	// Update allowed fields - only update if explicitly provided
	if payload.Pokok > 0 {
		existing.Pokok = payload.Pokok
	}
	// For bunga, check if it's explicitly set (even to 0) by checking if it's different from default
	// We should only update bunga if it's explicitly provided in the request
	if payload.Bunga > 0 {
		existing.Bunga = payload.Bunga
	}
	if payload.Denda >= 0 {
		existing.Denda = payload.Denda
	}
	if payload.TotalBayar > 0 {
		existing.TotalBayar = payload.TotalBayar
	}
	if payload.Status != "" {
		existing.Status = payload.Status
	}
	if !payload.TanggalBayar.IsZero() {
		existing.TanggalBayar = payload.TanggalBayar
	}

	// Recalculate total if components changed
	if payload.Pokok > 0 || payload.Bunga >= 0 || payload.Denda >= 0 {
		existing.TotalBayar = existing.Pokok + existing.Bunga + existing.Denda
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete removes an Angsuran
func (s *AngsuranService) Delete(requestorID uint, requestorRole string, id uint) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Super admin can delete any angsuran
	// Admin and member can only delete their own angsuran
	if requestorRole != "super_admin" && existing.UserID != requestorID {
		return errors.New("forbidden")
	}

	return s.repo.Delete(id)
}

// VerifyPayment allows admin to verify angsuran payment and update status
func (s *AngsuranService) VerifyPayment(requestorID uint, requestorRole string, id uint, status string) (*model.Angsuran, error) {
	// Only admin and super_admin can verify payments
	if requestorRole != "admin" && requestorRole != "super_admin" {
		return nil, errors.New("forbidden")
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check access permissions based on role
	if requestorRole == "super_admin" {
		// Super admin can verify any payment
	} else if requestorRole == "admin" {
		// Admin can verify payments for users they registered
		user, err := s.userRepo.FindByID(existing.UserID)
		if err != nil {
			return nil, err
		}
		if user.AdminID == nil || *user.AdminID != requestorID {
			return nil, errors.New("forbidden")
		}
	}

	// Validate status
	validStatuses := map[string]bool{"verified": true, "kurang": true, "lebih": true}
	if !validStatuses[status] {
		return nil, errors.New("invalid status for verification")
	}

	existing.Status = status

	// If payment is verified, update the pinjaman's remaining installments
	if status == "verified" {
		pinjaman, err := s.pinjamanRepo.GetByID(existing.PinjamanID)
		if err == nil && pinjaman.SisaAngsuran > 0 {
			pinjaman.SisaAngsuran--
			if pinjaman.SisaAngsuran == 0 {
				pinjaman.Status = "lunas"
			}
			s.pinjamanRepo.Update(pinjaman)
		}
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// GetPendingPayments returns angsuran with 'proses' status for admin verification
func (s *AngsuranService) GetPendingPayments(requestorID uint, requestorRole string) ([]model.Angsuran, error) {
	if requestorRole == "super_admin" {
		return s.repo.GetByStatus("proses", 0)
	}

	if requestorRole == "admin" {
		return s.repo.GetByStatus("proses", requestorID)
	}

	return nil, errors.New("forbidden")
}
