package repository

import (
	"backend-liburwan/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BackupAssignmentRepository struct {
	db *gorm.DB
}

func NewBackupAssignmentRepository(db *gorm.DB) *BackupAssignmentRepository {
	return &BackupAssignmentRepository{db: db}
}

func (r *BackupAssignmentRepository) DB() *gorm.DB {
	return r.db
}

func (r *BackupAssignmentRepository) GetByID(id uuid.UUID) (*model.BackupAssignment, error) {
	var backup model.BackupAssignment
	err := r.db.Preload("BackupKaryawan").Preload("Assigner").First(&backup, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &backup, nil
}

func (r *BackupAssignmentRepository) GetByJadwalLiburID(id uuid.UUID) (*model.BackupAssignment, error) {
	var backup model.BackupAssignment
	err := r.db.First(&backup, "jadwal_libur_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &backup, nil
}

func (r *BackupAssignmentRepository) Create(backup *model.BackupAssignment) error {
	return r.db.Create(backup).Error
}

func (r *BackupAssignmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.BackupAssignment{}, "id = ?", id).Error
}
