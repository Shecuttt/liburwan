package repository

import (
	"backend-liburwan/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokoRepository struct {
	db *gorm.DB
}

func NewTokoRepository(db *gorm.DB) *TokoRepository {
	return &TokoRepository{db: db}
}

func (r *TokoRepository) GetAll() ([]model.Toko, error) {
	var tokos []model.Toko
	err := r.db.Find(&tokos).Error
	return tokos, err
}

func (r *TokoRepository) GetByID(id uuid.UUID) (*model.Toko, error) {
	var toko model.Toko
	err := r.db.First(&toko, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &toko, nil
}

func (r *TokoRepository) Create(toko *model.Toko) error {
	return r.db.Create(toko).Error
}
