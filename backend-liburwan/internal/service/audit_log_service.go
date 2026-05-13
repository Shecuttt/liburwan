package service

import (
	"backend-liburwan/internal/model"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLogService struct {
}

func NewAuditLogService() *AuditLogService {
	return &AuditLogService{}
}

func (s *AuditLogService) Log(tx *gorm.DB, karyawanID *uuid.UUID, action string, entity string, entityID uuid.UUID, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	logEntry := &model.AuditLog{
		KaryawanID: karyawanID,
		Action:     action,
		Entity:     entity,
		EntityID:   entityID,
		Payload:    string(payloadBytes),
	}

	return tx.Create(logEntry).Error
}
