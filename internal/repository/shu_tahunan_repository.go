package repository

import (
	"koperasi-service/internal/model"

	"gorm.io/gorm"
)

// SHUTahunanRepository handles persistence for SHUTahunan entities
type SHUTahunanRepository struct {
	db *gorm.DB
}

// NewSHUTahunanRepository constructs a new repository instance
func NewSHUTahunanRepository(db *gorm.DB) *SHUTahunanRepository {
	return &SHUTahunanRepository{db: db}
}

// Create inserts a new SHUTahunan record
func (r *SHUTahunanRepository) Create(s *model.SHUTahunan) error {
	return r.db.Create(s).Error
}

// GetAll returns all SHU records
func (r *SHUTahunanRepository) GetAll() ([]model.SHUTahunan, error) {
	var list []model.SHUTahunan
	if err := r.db.Order("tahun DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetByID returns single SHU record by id
func (r *SHUTahunanRepository) GetByID(id uint) (*model.SHUTahunan, error) {
	var s model.SHUTahunan
	if err := r.db.First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// GetByTahun returns SHU record by year
func (r *SHUTahunanRepository) GetByTahun(tahun int) (*model.SHUTahunan, error) {
	var s model.SHUTahunan
	if err := r.db.Where("tahun = ?", tahun).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// Update persists changes to an existing SHUTahunan
func (r *SHUTahunanRepository) Update(s *model.SHUTahunan) error {
	return r.db.Save(s).Error
}

// Delete removes a SHUTahunan by id
func (r *SHUTahunanRepository) Delete(id uint) error {
	return r.db.Delete(&model.SHUTahunan{}, id).Error
}

// GetTotalSimpananByYear calculates total simpanan for a specific year
func (r *SHUTahunanRepository) GetTotalSimpananByYear(tahun int) (float64, error) {
	var total float64
	err := r.db.Model(&model.Simpanan{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("EXTRACT(YEAR FROM created_at) = ?", tahun).
		Scan(&total).Error
	return total, err
}

// GetSimpananByUserAndYear calculates simpanan per user for a specific year
func (r *SHUTahunanRepository) GetSimpananByUserAndYear(tahun int) (map[uint]float64, error) {
	type UserSimpanan struct {
		UserID uint    `json:"user_id"`
		Total  float64 `json:"total"`
	}

	var results []UserSimpanan
	err := r.db.Model(&model.Simpanan{}).
		Select("user_id, COALESCE(SUM(amount), 0) as total").
		Where("EXTRACT(YEAR FROM created_at) = ?", tahun).
		Group("user_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	userSimpanan := make(map[uint]float64)
	for _, result := range results {
		userSimpanan[result.UserID] = result.Total
	}

	return userSimpanan, nil
}

// GetTotalPenjualanByYear calculates total penjualan (loans) for a specific year
// Note: Using pinjaman as proxy for "penjualan" since it's the main transaction volume
func (r *SHUTahunanRepository) GetTotalPenjualanByYear(tahun int) (float64, error) {
	var total float64
	err := r.db.Model(&model.Pinjaman{}).
		Select("COALESCE(SUM(jumlah_pinjaman), 0)").
		Where("EXTRACT(YEAR FROM created_at) = ?", tahun).
		Scan(&total).Error
	return total, err
}

// GetPenjualanByUserAndYear calculates penjualan per user for a specific year
func (r *SHUTahunanRepository) GetPenjualanByUserAndYear(tahun int) (map[uint]float64, error) {
	type UserPenjualan struct {
		UserID uint    `json:"user_id"`
		Total  float64 `json:"total"`
	}

	var results []UserPenjualan
	err := r.db.Model(&model.Pinjaman{}).
		Select("user_id, COALESCE(SUM(jumlah_pinjaman), 0) as total").
		Where("EXTRACT(YEAR FROM created_at) = ?", tahun).
		Group("user_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	userPenjualan := make(map[uint]float64)
	for _, result := range results {
		userPenjualan[result.UserID] = result.Total
	}

	return userPenjualan, nil
}

// GetPendapatanOperasionalByYear calculates operational income for a specific year
// This includes income from loans (bunga), fees, and other operational activities
func (r *SHUTahunanRepository) GetPendapatanOperasionalByYear(tahun int) (float64, error) {
	var total float64

	// Calculate from loan interest (bunga from angsuran that are verified)
	err := r.db.Model(&model.Angsuran{}).
		Select("COALESCE(SUM(bunga), 0)").
		Where("EXTRACT(YEAR FROM created_at) = ? AND status = ?", tahun, "verified").
		Scan(&total).Error

	if err != nil {
		return 0, err
	}

	return total, nil
}

// GetPendapatanNonOperasionalByYear calculates non-operational income for a specific year
// This could include investment returns, grants, or other non-operational income
func (r *SHUTahunanRepository) GetPendapatanNonOperasionalByYear(tahun int) (float64, error) {
	// For now, this returns 0 as we don't have non-operational income tracking
	// This can be extended later when non-operational income sources are added
	return 0, nil
}

// GetAllUsers returns all users for SHU calculation
func (r *SHUTahunanRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
