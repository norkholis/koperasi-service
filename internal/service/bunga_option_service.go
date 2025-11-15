package service

import (
	"errors"
	"koperasi-service/internal/model"
	"koperasi-service/internal/repository"
)

type BungaOptionService interface {
	CreateBungaOption(userID uint, nama string, persen float64, deskripsi string) (*model.BungaOption, error)
	GetBungaOptionByID(id uint) (*model.BungaOption, error)
	GetAllBungaOptions() ([]model.BungaOption, error)
	GetActiveBungaOptions() ([]model.BungaOption, error)
	UpdateBungaOption(id uint, userID uint, nama string, persen float64, deskripsi string) error
	DeleteBungaOption(id uint, userID uint) error
	SetBungaOptionActive(id uint, userID uint, isActive bool) error
}

type bungaOptionService struct {
	bungaOptionRepo repository.BungaOptionRepository
	userRepo        *repository.UserRepository
}

func NewBungaOptionService(bungaOptionRepo repository.BungaOptionRepository, userRepo *repository.UserRepository) BungaOptionService {
	return &bungaOptionService{
		bungaOptionRepo: bungaOptionRepo,
		userRepo:        userRepo,
	}
}

func (s *bungaOptionService) CreateBungaOption(userID uint, nama string, persen float64, deskripsi string) (*model.BungaOption, error) {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return nil, errors.New("only admin can create bunga options")
	}

	// Validate percentage (should be positive)
	if persen <= 0 {
		return nil, errors.New("percentage must be greater than 0")
	}

	bungaOption := &model.BungaOption{
		Nama:      nama,
		Persen:    persen,
		Deskripsi: deskripsi,
		IsActive:  true,
		CreatedBy: userID,
	}

	err = s.bungaOptionRepo.Create(bungaOption)
	if err != nil {
		return nil, err
	}

	return bungaOption, nil
}

func (s *bungaOptionService) GetBungaOptionByID(id uint) (*model.BungaOption, error) {
	return s.bungaOptionRepo.GetByID(id)
}

func (s *bungaOptionService) GetAllBungaOptions() ([]model.BungaOption, error) {
	return s.bungaOptionRepo.GetAll()
}

func (s *bungaOptionService) GetActiveBungaOptions() ([]model.BungaOption, error) {
	return s.bungaOptionRepo.GetActiveOptions()
}

func (s *bungaOptionService) UpdateBungaOption(id uint, userID uint, nama string, persen float64, deskripsi string) error {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return errors.New("only admin can update bunga options")
	}

	// Validate percentage
	if persen <= 0 {
		return errors.New("percentage must be greater than 0")
	}

	bungaOption := &model.BungaOption{
		Nama:      nama,
		Persen:    persen,
		Deskripsi: deskripsi,
	}

	return s.bungaOptionRepo.Update(id, bungaOption)
}

func (s *bungaOptionService) DeleteBungaOption(id uint, userID uint) error {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return errors.New("only admin can delete bunga options")
	}

	return s.bungaOptionRepo.Delete(id)
}

func (s *bungaOptionService) SetBungaOptionActive(id uint, userID uint, isActive bool) error {
	// Check if user is admin or super_admin
	user, err := s.userRepo.FindByIDWithRole(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Role.Name != "admin" && user.Role.Name != "super_admin" {
		return errors.New("only admin can modify bunga option status")
	}

	return s.bungaOptionRepo.SetActive(id, isActive)
}
