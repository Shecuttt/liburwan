package service

import (
	"backend-liburwan/internal/model"
	"backend-liburwan/internal/repository"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrAlreadyAssigned = errors.New("ALREADY_ASSIGNED")
)

type BackupAssignmentService struct {
	repo         *repository.BackupAssignmentRepository
	jadwalRepo   *repository.JadwalLiburRepository
	auditService *AuditLogService
}

func NewBackupAssignmentService(repo *repository.BackupAssignmentRepository, jadwalRepo *repository.JadwalLiburRepository, auditService *AuditLogService) *BackupAssignmentService {
	return &BackupAssignmentService{repo: repo, jadwalRepo: jadwalRepo, auditService: auditService}
}

func (s *BackupAssignmentService) Create(jadwalLiburID, backupKaryawanID, assignedBy uuid.UUID) (*model.BackupAssignment, error) {
	jadwal, err := s.jadwalRepo.GetByID(jadwalLiburID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	existing, _ := s.repo.GetByJadwalLiburID(jadwalLiburID)
	if existing != nil {
		return nil, ErrAlreadyAssigned
	}

	if backupKaryawanID == jadwal.KaryawanID {
		return nil, errors.New("BACKUP_SELF")
	}

	onLeaveIDs, _ := s.jadwalRepo.GetEmployeesOnLeave(uuid.Nil, jadwal.Tanggal)
	for _, id := range onLeaveIDs {
		if id == backupKaryawanID {
			return nil, ErrBackupInvalid
		}
	}

	backup := &model.BackupAssignment{
		JadwalLiburID:    jadwalLiburID,
		BackupKaryawanID: backupKaryawanID,
		AssignedBy:       assignedBy,
	}

	err = s.repo.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(backup).Error; err != nil {
			return err
		}

		payload := map[string]interface{}{
			"backup_assignment": backup,
		}
		return s.auditService.Log(tx, &assignedBy, "CREATE_BACKUP_ASSIGNMENT", "backup_assignment", backup.ID, payload)
	})

	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(backup.ID)
}

func (s *BackupAssignmentService) Delete(id uuid.UUID, deletedBy uuid.UUID) error {
	backup, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		payload := map[string]interface{}{
			"backup_assignment": backup,
		}
		
		if err := s.auditService.Log(tx, &deletedBy, "DELETE_BACKUP_ASSIGNMENT", "backup_assignment", id, payload); err != nil {
			return err
		}

		return tx.Delete(&model.BackupAssignment{}, "id = ?", id).Error
	})
}
