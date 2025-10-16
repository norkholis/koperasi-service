package repository

import (
	"koperasi-service/internal/model"

	"gorm.io/gorm"
)

// SimpananRepository handles persistence for Simpanan entities.
type SimpananRepository struct {
	db *gorm.DB
}

// NewSimpananRepository constructs a new repository instance.
func NewSimpananRepository(db *gorm.DB) *SimpananRepository {
	return &SimpananRepository{db: db}
}

// Create inserts a new Simpanan record.
func (r *SimpananRepository) Create(s *model.Simpanan) error {
	return r.db.Create(s).Error
}

// GetAll returns all simpanan records; if userID > 0 it filters by user.
func (r *SimpananRepository) GetAll(userID uint) ([]model.Simpanan, error) {
	var list []model.Simpanan
	q := r.db
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetByID returns single simpanan by id.
func (r *SimpananRepository) GetByID(id uint) (*model.Simpanan, error) {
	var s model.Simpanan
	if err := r.db.First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// Update persists changes to an existing Simpanan.
func (r *SimpananRepository) Update(s *model.Simpanan) error {
	return r.db.Save(s).Error
}

// Delete removes a Simpanan by id.
func (r *SimpananRepository) Delete(id uint) error {
	return r.db.Delete(&model.Simpanan{}, id).Error
}
