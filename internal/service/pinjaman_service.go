package service

import (
	"errors"
	"fmt"
	"koperasi-service/internal/model"
	"koperasi-service/internal/repository"
	"time"
)

// PinjamanService handles business logic for Pinjaman with role constraints
type PinjamanService struct {
	repo            *repository.PinjamanRepository
	userRepo        *repository.UserRepository
	bungaOptionRepo repository.BungaOptionRepository
}

// NewPinjamanService creates a new service instance
func NewPinjamanService(repo *repository.PinjamanRepository, userRepo *repository.UserRepository, bungaOptionRepo repository.BungaOptionRepository) *PinjamanService {
	return &PinjamanService{repo: repo, userRepo: userRepo, bungaOptionRepo: bungaOptionRepo}
}

// Create adds a new Pinjaman (members can create for themselves, admins can create for any user)
func (s *PinjamanService) Create(requestorID uint, requestorRole string, p *model.Pinjaman) error {
	// Members can only create loans for themselves
	if requestorRole == "member" && p.UserID != requestorID {
		return errors.New("forbidden")
	}

	// If bunga_option_id is provided, get the interest rate from the option
	if p.BungaOptionID != nil {
		bungaOption, err := s.bungaOptionRepo.GetByID(*p.BungaOptionID)
		if err != nil {
			return errors.New("invalid bunga option")
		}
		if !bungaOption.IsActive {
			return errors.New("bunga option is not active")
		}
		// Copy the interest rate for historical record
		p.BungaPersen = bungaOption.Persen
	}

	// Validate that BungaPersen is set (either provided directly or copied from option)
	if p.BungaPersen == 0 {
		return errors.New("bunga_persen is required")
	}

	// Generate kode_pinjaman if not provided
	if p.KodePinjaman == "" {
		p.KodePinjaman = s.generateKodePinjaman()
	}

	// Set default values
	if p.TanggalPinjam.IsZero() {
		p.TanggalPinjam = time.Now()
	}
	if p.Status == "" {
		p.Status = "proses"
	}
	if p.SisaAngsuran == 0 {
		p.SisaAngsuran = p.LamaBulan
	}

	return s.repo.Create(p)
}

// List returns pinjaman list filtered by user unless role allows viewing all
func (s *PinjamanService) List(requestorID uint, requestorRole string) ([]model.Pinjaman, error) {
	if requestorRole == "super_admin" {
		// Super admin can see all loans
		return s.repo.GetAll(0)
	}

	if requestorRole == "admin" {
		// Admin can see loans from users they registered AND their own loans
		adminLoans, err := s.repo.GetByAdminUserID(requestorID)
		if err != nil {
			return nil, err
		}

		// Also get admin's own loans
		ownLoans, err := s.repo.GetAll(requestorID)
		if err != nil {
			return nil, err
		}

		// Combine both lists (avoiding duplicates)
		allLoans := make([]model.Pinjaman, 0)
		allLoans = append(allLoans, adminLoans...)

		// Add own loans if not already included
		for _, ownLoan := range ownLoans {
			duplicate := false
			for _, adminLoan := range adminLoans {
				if adminLoan.ID == ownLoan.ID {
					duplicate = true
					break
				}
			}
			if !duplicate {
				allLoans = append(allLoans, ownLoan)
			}
		}

		return allLoans, nil
	}

	// Members can only see their own loans
	return s.repo.GetAll(requestorID)
}

// Get returns Pinjaman by id with access control
func (s *PinjamanService) Get(requestorID uint, requestorRole string, id uint) (*model.Pinjaman, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Super admin can view any loan
	if requestorRole == "super_admin" {
		return p, nil
	}

	// Admin can view loans from users they registered or their own loans
	if requestorRole == "admin" {
		// Check if loan owner is their own or user they registered
		user, err := s.userRepo.FindByID(p.UserID)
		if err != nil {
			return nil, err
		}

		// Allow if it's admin's own loan or if user was registered by this admin
		if p.UserID == requestorID || (user.AdminID != nil && *user.AdminID == requestorID) {
			return p, nil
		}
		return nil, errors.New("forbidden")
	}

	// Members can only view their own loans
	if p.UserID != requestorID {
		return nil, errors.New("forbidden")
	}

	return p, nil
}

// Update modifies an existing Pinjaman after verifying access
func (s *PinjamanService) Update(requestorID uint, requestorRole string, id uint, payload *model.Pinjaman) (*model.Pinjaman, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check access permissions based on role
	if requestorRole == "super_admin" {
		// Super admin can update any loan
	} else if requestorRole == "admin" {
		// Admin can update loans from users they registered
		user, err := s.userRepo.FindByID(existing.UserID)
		if err != nil {
			return nil, err
		}
		if user.AdminID == nil || *user.AdminID != requestorID {
			return nil, errors.New("forbidden")
		}
	} else {
		// Members can only update their own loans
		if existing.UserID != requestorID {
			return nil, errors.New("forbidden")
		}
	}

	// Update allowed fields
	if payload.JumlahPinjaman > 0 {
		existing.JumlahPinjaman = payload.JumlahPinjaman
	}
	// Only update BungaPersen if explicitly provided (> 0)
	// This prevents overwriting existing interest rate with default 0 value
	if payload.BungaPersen > 0 {
		existing.BungaPersen = payload.BungaPersen
	}
	if payload.LamaBulan > 0 {
		existing.LamaBulan = payload.LamaBulan
		// If loan duration changes, update remaining installments accordingly
		// Only if loan is still in "proses" status
		if existing.Status == "proses" {
			existing.SisaAngsuran = payload.LamaBulan
		}
	}
	if payload.JumlahAngsuran > 0 {
		existing.JumlahAngsuran = payload.JumlahAngsuran
	}
	// Note: SisaAngsuran cannot be directly updated via API
	// It's only decremented by the system when payments are verified
	if payload.Status != "" {
		existing.Status = payload.Status
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete removes a Pinjaman after verifying access
func (s *PinjamanService) Delete(requestorID uint, requestorRole string, id uint) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Check access permissions based on role
	if requestorRole == "super_admin" {
		// Super admin can delete any loan
	} else if requestorRole == "admin" {
		// Admin can delete loans from users they registered
		user, err := s.userRepo.FindByID(existing.UserID)
		if err != nil {
			return err
		}
		if user.AdminID == nil || *user.AdminID != requestorID {
			return errors.New("forbidden")
		}
	} else {
		// Members can only delete their own loans
		if existing.UserID != requestorID {
			return errors.New("forbidden")
		}
	}

	return s.repo.Delete(id)
}

// Approve approves a pinjaman (admin can approve loans from their registered users)
func (s *PinjamanService) Approve(requestorID uint, requestorRole string, pinjamanID uint, reason string) error {
	// Get the pinjaman
	pinjaman, err := s.repo.GetByID(pinjamanID)
	if err != nil {
		return errors.New("pinjaman not found")
	}

	// Check if pinjaman is in "proses" status
	if pinjaman.Status != "proses" {
		return errors.New("invalid pinjaman status")
	}

	// Check permissions
	if requestorRole == "super_admin" {
		// Super admin can approve any loan
	} else if requestorRole == "admin" {
		// Admin can only approve loans from users they registered
		user, err := s.userRepo.FindByID(pinjaman.UserID)
		if err != nil {
			return err
		}
		if user.AdminID == nil || *user.AdminID != requestorID {
			return errors.New("forbidden")
		}
	} else {
		// Members cannot approve loans
		return errors.New("forbidden")
	}

	// Update status to approved
	pinjaman.Status = "disetujui"

	return s.repo.Update(pinjaman)
}

// Reject rejects a pinjaman (admin can reject loans from their registered users)
func (s *PinjamanService) Reject(requestorID uint, requestorRole string, pinjamanID uint, reason string) error {
	// Get the pinjaman
	pinjaman, err := s.repo.GetByID(pinjamanID)
	if err != nil {
		return errors.New("pinjaman not found")
	}

	// Check if pinjaman is in "proses" status
	if pinjaman.Status != "proses" {
		return errors.New("invalid pinjaman status")
	}

	// Check permissions
	if requestorRole == "super_admin" {
		// Super admin can reject any loan
	} else if requestorRole == "admin" {
		// Admin can only reject loans from users they registered
		user, err := s.userRepo.FindByID(pinjaman.UserID)
		if err != nil {
			return err
		}
		if user.AdminID == nil || *user.AdminID != requestorID {
			return errors.New("forbidden")
		}
	} else {
		// Members cannot reject loans
		return errors.New("forbidden")
	}

	// Update status to rejected
	pinjaman.Status = "ditolak"

	return s.repo.Update(pinjaman)
}

// generateKodePinjaman creates a unique loan code
func (s *PinjamanService) generateKodePinjaman() string {
	return fmt.Sprintf("PJM%d", time.Now().Unix())
}
