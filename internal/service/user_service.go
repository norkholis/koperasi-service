package service

import (
	"errors"
	"koperasi-service/internal/model"
	"koperasi-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// UserService handles user CRUD with role constraints.
type UserService struct {
	repo         *repository.UserRepository
	simpananRepo *repository.SimpananRepository
}

// NewUserService constructs a new UserService.
func NewUserService(repo *repository.UserRepository, simpananRepo *repository.SimpananRepository) *UserService {
	return &UserService{repo: repo, simpananRepo: simpananRepo}
}

// ListUsers returns all users; only super_admin can list all, admin can list their registered users.
func (s *UserService) ListUsers(requestorID uint, requestorRole string) ([]model.User, error) {
	if requestorRole == "super_admin" {
		// Super admin can see all users
		return s.repo.List(true, 0)
	} else if requestorRole == "admin" {
		// Admin can only see users they registered
		return s.repo.List(true, requestorID)
	}
	return nil, errors.New("forbidden")
}

// GetUser returns a user by id; super_admin can see any; admin can see their registered users; others only themselves.
func (s *UserService) GetUser(requestorID uint, requestorRole string, targetID uint) (*model.User, error) {
	user, err := s.repo.FindByIDWithRole(targetID)
	if err != nil {
		return nil, err
	}

	if requestorRole == "super_admin" {
		return user, nil
	}

	if requestorRole == "admin" {
		// Admin can see users they registered or themselves
		if user.AdminID != nil && *user.AdminID == requestorID || user.ID == requestorID {
			return user, nil
		}
		return nil, errors.New("forbidden")
	}

	// Members can only see themselves
	if requestorID != targetID {
		return nil, errors.New("forbidden")
	}

	return user, nil
}

// CreateUser creates a new user; super_admin and admin can create users.
func (s *UserService) CreateUser(requestorID uint, requestorRole string, u *model.User) error {
	if requestorRole != "super_admin" && requestorRole != "admin" {
		return errors.New("forbidden")
	}

	// Only set admin_id for admin users, not for super_admin
	if requestorRole == "admin" {
		u.AdminID = &requestorID
	}
	// Super admin doesn't set admin_id (can manage all users)

	if u.Password != "" {
		h, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(h)
	}

	// Create user first
	if err := s.repo.Create(u); err != nil {
		return err
	}

	// Initialize user wallets (3 types)
	if err := s.simpananRepo.InitializeUserWallets(u.ID); err != nil {
		// If wallet initialization fails, we should rollback user creation
		// For now, we'll just return the error
		return err
	}

	return nil
}

// UpdateUser updates target user; super_admin any; admin their registered users; others only themselves. Role changes only by super_admin.
func (s *UserService) UpdateUser(requestorID uint, requestorRole string, targetID uint, email, name, address, phoneNumber, nik string, password *string, roleID *uint) (*model.User, error) {
	u, err := s.repo.FindByID(targetID)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if requestorRole == "super_admin" {
		// Super admin can update any user (no restrictions)
	} else if requestorRole == "admin" {
		// Admin can only update users they registered or themselves
		if u.AdminID == nil || *u.AdminID != requestorID {
			if u.ID != requestorID {
				return nil, errors.New("forbidden")
			}
		}
	} else {
		// Members can only update themselves
		if requestorID != targetID {
			return nil, errors.New("forbidden")
		}
	}

	if email != "" {
		u.Email = email
	}
	if name != "" {
		u.Name = name
	}
	if address != "" {
		u.Address = address
	}
	if phoneNumber != "" {
		u.PhoneNumber = phoneNumber
	}
	if nik != "" {
		u.NIK = nik
	}
	if password != nil && *password != "" {
		h, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		u.Password = string(h)
	}
	// Only super_admin can change roles
	if roleID != nil && requestorRole == "super_admin" {
		u.RoleID = *roleID
	}
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

// DeleteUser deletes target user; super_admin any; admin their registered users; member only themselves
func (s *UserService) DeleteUser(requestorID uint, requestorRole string, targetID uint) error {
	u, err := s.repo.FindByID(targetID)
	if err != nil {
		return err
	}

	// Check permissions
	if requestorRole == "super_admin" {
		// Super admin can delete any user (no restrictions)
	} else if requestorRole == "admin" {
		// Admin can only delete users they registered or themselves
		if u.AdminID == nil || *u.AdminID != requestorID {
			if u.ID != requestorID {
				return errors.New("forbidden")
			}
		}
	} else {
		// Members can only delete themselves
		if requestorID != targetID {
			return errors.New("forbidden")
		}
	}

	return s.repo.Delete(targetID)
}
