package repository

import (
	"koperasi-service/internal/model"

	"gorm.io/gorm"
)

// SHUAnggotaRepository handles database operations for SHU Anggota records
type SHUAnggotaRepository struct {
	db *gorm.DB
}

// NewSHUAnggotaRepository creates a new repository instance
func NewSHUAnggotaRepository(db *gorm.DB) *SHUAnggotaRepository {
	return &SHUAnggotaRepository{db: db}
}

// Create saves a new SHU Anggota record
func (r *SHUAnggotaRepository) Create(shuAnggota *model.SHUAnggotaRecord) error {
	return r.db.Create(shuAnggota).Error
}

// GetByID retrieves a SHU Anggota record by ID
func (r *SHUAnggotaRepository) GetByID(id uint) (*model.SHUAnggotaRecord, error) {
	var shuAnggota model.SHUAnggotaRecord
	err := r.db.Preload("SHU").Preload("User").First(&shuAnggota, id).Error
	return &shuAnggota, err
}

// GetBySHUIDAndUserID retrieves a SHU Anggota record by SHU ID and User ID
func (r *SHUAnggotaRepository) GetBySHUIDAndUserID(shuID, userID uint) (*model.SHUAnggotaRecord, error) {
	var shuAnggota model.SHUAnggotaRecord
	err := r.db.Preload("SHU").Preload("User").
		Where("id_shu = ? AND id_anggota = ?", shuID, userID).
		First(&shuAnggota).Error
	return &shuAnggota, err
}

// GetBySHUID retrieves all SHU Anggota records for a specific SHU
func (r *SHUAnggotaRepository) GetBySHUID(shuID uint) ([]model.SHUAnggotaRecord, error) {
	var shuAnggotas []model.SHUAnggotaRecord
	err := r.db.Preload("SHU").Preload("User").
		Where("id_shu = ?", shuID).
		Find(&shuAnggotas).Error
	return shuAnggotas, err
}

// GetByUserID retrieves all SHU Anggota records for a specific user
func (r *SHUAnggotaRepository) GetByUserID(userID uint) ([]model.SHUAnggotaRecord, error) {
	var shuAnggotas []model.SHUAnggotaRecord
	err := r.db.Preload("SHU").Preload("User").
		Where("id_anggota = ?", userID).
		Find(&shuAnggotas).Error
	return shuAnggotas, err
}

// Update modifies an existing SHU Anggota record
func (r *SHUAnggotaRepository) Update(id uint, shuAnggota *model.SHUAnggotaRecord) error {
	return r.db.Where("id_shu_anggota = ?", id).Updates(shuAnggota).Error
}

// Delete removes a SHU Anggota record
func (r *SHUAnggotaRepository) Delete(id uint) error {
	return r.db.Delete(&model.SHUAnggotaRecord{}, id).Error
}

// List retrieves all SHU Anggota records with pagination
func (r *SHUAnggotaRepository) List(offset, limit int) ([]model.SHUAnggotaRecord, error) {
	var shuAnggotas []model.SHUAnggotaRecord
	err := r.db.Preload("SHU").Preload("User").
		Offset(offset).Limit(limit).
		Find(&shuAnggotas).Error
	return shuAnggotas, err
}

// CheckExists checks if a SHU Anggota record already exists for a user and SHU
func (r *SHUAnggotaRepository) CheckExists(shuID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.SHUAnggotaRecord{}).
		Where("id_shu = ? AND id_anggota = ?", shuID, userID).
		Count(&count).Error
	return count > 0, err
}
