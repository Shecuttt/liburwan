package repository

import (
	"backend-liburwan/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KaryawanRepository struct {
	db *gorm.DB
}

func NewKaryawanRepository(db *gorm.DB) *KaryawanRepository {
	return &KaryawanRepository{db: db}
}

func (r *KaryawanRepository) DB() *gorm.DB {
	return r.db
}

func (r *KaryawanRepository) GetAll(tokoID string) ([]model.Karyawan, error) {
	var karyawans []model.Karyawan
	query := r.db.Preload("Toko")
	if tokoID != "" {
		query = query.Where("toko_id = ?", tokoID)
	}
	err := query.Find(&karyawans).Error
	return karyawans, err
}

func (r *KaryawanRepository) CountInStore(tokoID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.Karyawan{}).Where("toko_id = ?", tokoID).Count(&count).Error
	return count, err
}

func (r *KaryawanRepository) GetByEmail(email string) (*model.Karyawan, error) {
	var karyawan model.Karyawan
	err := r.db.Preload("Toko").Where("email = ?", email).First(&karyawan).Error
	if err != nil {
		return nil, err
	}
	return &karyawan, nil
}

func (r *KaryawanRepository) UpdateGoogleID(id uuid.UUID, googleID string) error {
	return r.db.Model(&model.Karyawan{}).Where("id = ?", id).Update("google_id", googleID).Error
}

func (r *KaryawanRepository) GetByID(id uuid.UUID) (*model.Karyawan, error) {
	var karyawan model.Karyawan
	err := r.db.Preload("Toko").First(&karyawan, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &karyawan, nil
}

func (r *KaryawanRepository) Create(karyawan *model.Karyawan) error {
	return r.db.Create(karyawan).Error
}

func (r *KaryawanRepository) Update(id uuid.UUID, data map[string]interface{}) error {
	return r.db.Model(&model.Karyawan{}).Where("id = ?", id).Updates(data).Error
}
