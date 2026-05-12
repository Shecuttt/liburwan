package service

import (
	"backend-liburwan/internal/model"
	"backend-liburwan/internal/repository"
	"github.com/google/uuid"
)

type TokoService struct {
	repo *repository.TokoRepository
}

func NewTokoService(repo *repository.TokoRepository) *TokoService {
	return &TokoService{repo: repo}
}

func (s *TokoService) GetAll() ([]model.Toko, error) {
	return s.repo.GetAll()
}

func (s *TokoService) GetByID(id uuid.UUID) (*model.Toko, error) {
	return s.repo.GetByID(id)
}

func (s *TokoService) Create(toko *model.Toko) error {
	return s.repo.Create(toko)
}
