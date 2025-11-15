package repository

import (
	"koperasi-service/internal/model"

	"gorm.io/gorm"
)

type BungaOptionRepository interface {
	Create(bungaOption *model.BungaOption) error
	GetByID(id uint) (*model.BungaOption, error)
	GetAll() ([]model.BungaOption, error)
	GetActiveOptions() ([]model.BungaOption, error)
	Update(id uint, bungaOption *model.BungaOption) error
	Delete(id uint) error
	SetActive(id uint, isActive bool) error
}

type bungaOptionRepository struct {
	db *gorm.DB
}

func NewBungaOptionRepository(db *gorm.DB) BungaOptionRepository {
	return &bungaOptionRepository{db: db}
}

func (r *bungaOptionRepository) Create(bungaOption *model.BungaOption) error {
	return r.db.Create(bungaOption).Error
}

func (r *bungaOptionRepository) GetByID(id uint) (*model.BungaOption, error) {
	var bungaOption model.BungaOption
	err := r.db.Preload("CreatedByUser").First(&bungaOption, id).Error
	if err != nil {
		return nil, err
	}
	return &bungaOption, nil
}

func (r *bungaOptionRepository) GetAll() ([]model.BungaOption, error) {
	var bungaOptions []model.BungaOption
	err := r.db.Preload("CreatedByUser").Find(&bungaOptions).Error
	return bungaOptions, err
}

func (r *bungaOptionRepository) GetActiveOptions() ([]model.BungaOption, error) {
	var bungaOptions []model.BungaOption
	err := r.db.Where("is_active = ?", true).Preload("CreatedByUser").Find(&bungaOptions).Error
	return bungaOptions, err
}

func (r *bungaOptionRepository) Update(id uint, bungaOption *model.BungaOption) error {
	return r.db.Model(&model.BungaOption{}).Where("id = ?", id).Updates(bungaOption).Error
}

func (r *bungaOptionRepository) Delete(id uint) error {
	return r.db.Delete(&model.BungaOption{}, id).Error
}

func (r *bungaOptionRepository) SetActive(id uint, isActive bool) error {
	return r.db.Model(&model.BungaOption{}).Where("id = ?", id).Update("is_active", isActive).Error
}
