package service

import (
	"errors"
	"koperasi-service/internal/model"
	"koperasi-service/internal/repository"
)

// SimpananService contains business logic for Simpanan.
type SimpananService struct {
	repo *repository.SimpananRepository
}

// NewSimpananService creates a new service instance.
func NewSimpananService(repo *repository.SimpananRepository) *SimpananService {
	return &SimpananService{repo: repo}
}

// Create adds a new Simpanan.
func (s *SimpananService) Create(sm *model.Simpanan) error {
	return s.repo.Create(sm)
}

// List returns simpanan list filtered by user unless allowAll is true.
func (s *SimpananService) List(userID uint, allowAll bool) ([]model.Simpanan, error) {
	// if allowAll, ignore userID filter by passing 0
	if allowAll {
		return s.repo.GetAll(0)
	}
	return s.repo.GetAll(userID)
}

// Get returns Simpanan by id.
func (s *SimpananService) Get(id uint) (*model.Simpanan, error) {
	return s.repo.GetByID(id)
}

// Update modifies an existing Simpanan after verifying ownership or admin rights.
func (s *SimpananService) Update(id uint, userID uint, allowAll bool, payload *model.Simpanan) (*model.Simpanan, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if !allowAll && existing.UserID != userID {
		return nil, errors.New("forbidden")
	}
	existing.Type = payload.Type
	existing.Amount = payload.Amount
	existing.Description = payload.Description
	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// Delete removes a Simpanan after verifying ownership or admin rights.
func (s *SimpananService) Delete(id uint, userID uint, allowAll bool) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if !allowAll && existing.UserID != userID {
		return errors.New("forbidden")
	}
	return s.repo.Delete(id)
}
