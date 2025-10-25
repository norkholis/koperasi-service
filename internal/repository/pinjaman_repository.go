package repository

import (
	"koperasi-service/internal/model"

	"gorm.io/gorm"
)

// PinjamanRepository handles persistence for Pinjaman entities
type PinjamanRepository struct {
	db *gorm.DB
}

// NewPinjamanRepository constructs a new repository instance
func NewPinjamanRepository(db *gorm.DB) *PinjamanRepository {
	return &PinjamanRepository{db: db}
}

// Create inserts a new Pinjaman record
func (r *PinjamanRepository) Create(p *model.Pinjaman) error {
	return r.db.Create(p).Error
}

// GetAll returns all pinjaman records; if userID > 0 it filters by user
func (r *PinjamanRepository) GetAll(userID uint) ([]model.Pinjaman, error) {
	var list []model.Pinjaman
	q := r.db.Preload("User")
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetByID returns single pinjaman by id with user preloaded
func (r *PinjamanRepository) GetByID(id uint) (*model.Pinjaman, error) {
	var p model.Pinjaman
	if err := r.db.Preload("User").First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// Update persists changes to an existing Pinjaman
func (r *PinjamanRepository) Update(p *model.Pinjaman) error {
	return r.db.Save(p).Error
}

// Delete removes a Pinjaman by id
func (r *PinjamanRepository) Delete(id uint) error {
	return r.db.Delete(&model.Pinjaman{}, id).Error
}

// GetByKodePinjaman finds pinjaman by kode_pinjaman
func (r *PinjamanRepository) GetByKodePinjaman(kode string) (*model.Pinjaman, error) {
	var p model.Pinjaman
	if err := r.db.Preload("User").Where("kode_pinjaman = ?", kode).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// GetByAdminUserID returns pinjaman records for users managed by admin
func (r *PinjamanRepository) GetByAdminUserID(adminID uint) ([]model.Pinjaman, error) {
	var list []model.Pinjaman
	if err := r.db.Preload("User").
		Joins("JOIN users ON pinjaman.user_id = users.id").
		Where("users.admin_id = ?", adminID).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
