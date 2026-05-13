package service

import (
	"backend-liburwan/internal/model"
	"backend-liburwan/internal/repository"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KonfigurasiService struct {
	repo         *repository.KonfigurasiRepository
	auditService *AuditLogService
}

func NewKonfigurasiService(repo *repository.KonfigurasiRepository, auditService *AuditLogService) *KonfigurasiService {
	return &KonfigurasiService{repo: repo, auditService: auditService}
}

func (s *KonfigurasiService) GetAll() ([]model.Konfigurasi, error) {
	return s.repo.GetAll()
}

func (s *KonfigurasiService) Update(key string, value string, adminID uuid.UUID) error {
	var oldConfig model.Konfigurasi
	err := s.repo.DB().Where("key = ?", key).First(&oldConfig).Error
	if err != nil {
		return err // Not found or DB error
	}

	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Konfigurasi{}).Where("key = ?", key).Update("value", value).Error; err != nil {
			return err
		}

		payload := map[string]interface{}{
			"before": map[string]interface{}{
				"value": oldConfig.Value,
			},
			"after": map[string]interface{}{
				"value": value,
			},
		}

		return s.auditService.Log(tx, &adminID, "UPDATE_KONFIGURASI", "konfigurasi", oldConfig.ID, payload)
	})
}

func (s *KonfigurasiService) GetConfigMap() (map[string]string, error) {
	configs, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, c := range configs {
		m[c.Key] = c.Value
	}
	return m, nil
}

func (s *KonfigurasiService) GetInt(m map[string]string, key string, fallback int) int {
	if v, ok := m[key]; ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
