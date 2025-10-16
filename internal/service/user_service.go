package service

import (
	"errors"
	"koperasi-service/internal/model"
	"koperasi-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// UserService handles user CRUD with role constraints.
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService constructs a new UserService.
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// ListUsers returns all users; only super_admin is permitted.
func (s *UserService) ListUsers(requestorRole string) ([]model.User, error) {
	if requestorRole != "super_admin" {
		return nil, errors.New("forbidden")
	}
	return s.repo.List(true)
}

// GetUser returns a user by id; super_admin can view any; others only themselves.
func (s *UserService) GetUser(requestorID uint, requestorRole string, targetID uint) (*model.User, error) {
	if requestorRole != "super_admin" && requestorID != targetID {
		return nil, errors.New("forbidden")
	}
	return s.repo.FindByIDWithRole(targetID)
}

// CreateUser creates a new user; restricted to super_admin.
func (s *UserService) CreateUser(requestorRole string, u *model.User) error {
	if requestorRole != "super_admin" {
		return errors.New("forbidden")
	}
	if u.Password != "" {
		h, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(h)
	}
	return s.repo.Create(u)
}

// UpdateUser updates an existing user; super_admin any user; others only themselves. Role changes only by super_admin.
func (s *UserService) UpdateUser(requestorID uint, requestorRole string, targetID uint, email string, password *string, roleID *uint) (*model.User, error) {
	u, err := s.repo.FindByID(targetID)
	if err != nil {
		return nil, err
	}
	if requestorRole != "super_admin" && requestorID != targetID {
		return nil, errors.New("forbidden")
	}
	if email != "" {
		u.Email = email
	}
	if password != nil && *password != "" {
		h, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		u.Password = string(h)
	}
	if roleID != nil && requestorRole == "super_admin" {
		u.RoleID = *roleID
	}
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

// DeleteUser deletes a user; super_admin any user; others only themselves.
func (s *UserService) DeleteUser(requestorID uint, requestorRole string, targetID uint) error {
	if requestorRole != "super_admin" && requestorID != targetID {
		return errors.New("forbidden")
	}
	return s.repo.Delete(targetID)
}
