package service

import (
	"backend-liburwan/internal/model"
	"backend-liburwan/internal/repository"
	"github.com/google/uuid"
)

type KaryawanService struct {
	repo *repository.KaryawanRepository
}

func NewKaryawanService(repo *repository.KaryawanRepository) *KaryawanService {
	return &KaryawanService{repo: repo}
}

func (s *KaryawanService) GetAll(tokoID string) ([]model.Karyawan, error) {
	return s.repo.GetAll(tokoID)
}

func (s *KaryawanService) GetByID(id uuid.UUID) (*model.Karyawan, error) {
	return s.repo.GetByID(id)
}

func (s *KaryawanService) Create(karyawan *model.Karyawan) error {
	return s.repo.Create(karyawan)
}

func (s *KaryawanService) Update(id uuid.UUID, data map[string]interface{}) error {
	return s.repo.Update(id, data)
}
