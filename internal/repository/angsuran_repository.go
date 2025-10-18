package repository

import (
	"koperasi-service/internal/model"

	"gorm.io/gorm"
)

// AngsuranRepository handles persistence for Angsuran entities
type AngsuranRepository struct {
	db *gorm.DB
}

// NewAngsuranRepository constructs a new repository instance
func NewAngsuranRepository(db *gorm.DB) *AngsuranRepository {
	return &AngsuranRepository{db: db}
}

// Create inserts a new Angsuran record
func (r *AngsuranRepository) Create(a *model.Angsuran) error {
	return r.db.Create(a).Error
}

// GetAll returns all angsuran records; filters by userID and/or pinjamanID if provided
func (r *AngsuranRepository) GetAll(userID uint, pinjamanID uint) ([]model.Angsuran, error) {
	var list []model.Angsuran
	q := r.db.Preload("Pinjaman").Preload("User")
	
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	if pinjamanID > 0 {
		q = q.Where("pinjaman_id = ?", pinjamanID)
	}
	
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetByID returns single angsuran by id with relations preloaded
func (r *AngsuranRepository) GetByID(id uint) (*model.Angsuran, error) {
	var a model.Angsuran
	if err := r.db.Preload("Pinjaman").Preload("User").First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

// Update persists changes to an existing Angsuran
func (r *AngsuranRepository) Update(a *model.Angsuran) error {
	return r.db.Save(a).Error
}

// Delete removes an Angsuran by id
func (r *AngsuranRepository) Delete(id uint) error {
	return r.db.Delete(&model.Angsuran{}, id).Error
}

// GetByPinjamanAndAngsuranKe finds angsuran by pinjaman ID and sequence number
func (r *AngsuranRepository) GetByPinjamanAndAngsuranKe(pinjamanID uint, angsuranKe int) (*model.Angsuran, error) {
	var a model.Angsuran
	if err := r.db.Preload("Pinjaman").Preload("User").
		Where("pinjaman_id = ? AND angsuran_ke = ?", pinjamanID, angsuranKe).
		First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

// GetByStatus returns angsuran records filtered by status
func (r *AngsuranRepository) GetByStatus(status string, userID uint) ([]model.Angsuran, error) {
	var list []model.Angsuran
	q := r.db.Preload("Pinjaman").Preload("User").Where("status = ?", status)
	
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}