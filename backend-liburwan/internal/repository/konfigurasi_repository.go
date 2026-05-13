package repository

import (
	"backend-liburwan/internal/model"
	"gorm.io/gorm"
)

type KonfigurasiRepository struct {
	db *gorm.DB
}

func NewKonfigurasiRepository(db *gorm.DB) *KonfigurasiRepository {
	return &KonfigurasiRepository{db: db}
}

func (r *KonfigurasiRepository) DB() *gorm.DB {
	return r.db
}

func (r *KonfigurasiRepository) GetAll() ([]model.Konfigurasi, error) {
	var configs []model.Konfigurasi
	err := r.db.Find(&configs).Error
	return configs, err
}

func (r *KonfigurasiRepository) GetByKey(key string) (*model.Konfigurasi, error) {
	var config model.Konfigurasi
	err := r.db.Where("key = ?", key).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *KonfigurasiRepository) Update(key string, value string) error {
	return r.db.Model(&model.Konfigurasi{}).Where("key = ?", key).Update("value", value).Error
}
