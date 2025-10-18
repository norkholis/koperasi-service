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
	repo *repository.PinjamanRepository
}

// NewPinjamanService creates a new service instance
func NewPinjamanService(repo *repository.PinjamanRepository) *PinjamanService {
	return &PinjamanService{repo: repo}
}

// Create adds a new Pinjaman (members can create for themselves, admins can create for any user)
func (s *PinjamanService) Create(requestorID uint, requestorRole string, p *model.Pinjaman) error {
	// Members can only create loans for themselves
	if requestorRole == "member" && p.UserID != requestorID {
		return errors.New("forbidden")
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

	// Admin and member can only see their own loans
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

	// Admin and member can only view their own loans
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

	// Super admin can update any loan
	// Admin and member can only update their own loans
	if requestorRole != "super_admin" && existing.UserID != requestorID {
		return nil, errors.New("forbidden")
	}

	// Update allowed fields
	if payload.JumlahPinjaman > 0 {
		existing.JumlahPinjaman = payload.JumlahPinjaman
	}
	if payload.BungaPersen >= 0 {
		existing.BungaPersen = payload.BungaPersen
	}
	if payload.LamaBulan > 0 {
		existing.LamaBulan = payload.LamaBulan
	}
	if payload.JumlahAngsuran > 0 {
		existing.JumlahAngsuran = payload.JumlahAngsuran
	}
	if payload.SisaAngsuran >= 0 {
		existing.SisaAngsuran = payload.SisaAngsuran
	}
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

	// Super admin can delete any loan
	// Admin and member can only delete their own loans
	if requestorRole != "super_admin" && existing.UserID != requestorID {
		return errors.New("forbidden")
	}

	return s.repo.Delete(id)
}

// generateKodePinjaman creates a unique loan code
func (s *PinjamanService) generateKodePinjaman() string {
	return fmt.Sprintf("PJM%d", time.Now().Unix())
}
