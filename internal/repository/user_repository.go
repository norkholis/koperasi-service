package repository

import (
	"koperasi-service/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByIDWithRole(id uint) (*model.User, error) {
	var u model.User
	err := r.db.Preload("Role").First(&u, id).Error
	return &u, err
}

// List returns all users (with role preload if withRole).
func (r *UserRepository) List(withRole bool, adminID uint) ([]model.User, error) {
	var users []model.User
	q := r.db
	if withRole {
		q = q.Preload("Role")
	}
	if adminID > 0 {
		q = q.Where("admin_id = ?", adminID)
	}
	if err := q.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// FindByID returns user by id (no preload).
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var u model.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Update saves user changes.
func (r *UserRepository) Update(u *model.User) error {
	return r.db.Save(u).Error
}

// Delete removes user by id.
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
