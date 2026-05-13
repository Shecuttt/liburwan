package service

import (
	"backend-liburwan/internal/lib/timeutil"
	"backend-liburwan/internal/model"
	"backend-liburwan/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrOutOfWindow       = errors.New("OUT_OF_WINDOW")
	ErrKuotaHabis        = errors.New("KUOTA_HABIS")
	ErrBackupRequired    = errors.New("BACKUP_REQUIRED")
	ErrBackupInvalid     = errors.New("BACKUP_INVALID")
	ErrBackupSelf        = errors.New("BACKUP_SELF")
	ErrNoBackupAvailable = errors.New("NO_BACKUP_AVAILABLE")
	ErrTanggalTerlewat   = errors.New("TANGGAL_TERLEWAT")
	ErrNotFound          = errors.New("NOT_FOUND")
	ErrUnauthorized      = errors.New("UNAUTHORIZED_ACTION")
)

type JadwalLiburService struct {
	repo          *repository.JadwalLiburRepository
	karyawanRepo  *repository.KaryawanRepository
	configService *KonfigurasiService
	auditService  *AuditLogService
}

func NewJadwalLiburService(repo *repository.JadwalLiburRepository, karyawanRepo *repository.KaryawanRepository, configService *KonfigurasiService, auditService *AuditLogService) *JadwalLiburService {
	return &JadwalLiburService{repo: repo, karyawanRepo: karyawanRepo, configService: configService, auditService: auditService}
}

func (s *JadwalLiburService) GetAll(karyawanID, tokoID, bulan string) ([]model.JadwalLibur, error) {
	return s.repo.GetAll(karyawanID, tokoID, bulan)
}

func (s *JadwalLiburService) GetByID(id uuid.UUID) (*model.JadwalLibur, error) {
	return s.repo.GetByID(id)
}

func (s *JadwalLiburService) CheckAvailability(karyawanID uuid.UUID, tanggal time.Time) (int, bool, []model.Karyawan, error) {
	configs, _ := s.configService.GetConfigMap()
	maxLibur := s.configService.GetInt(configs, "maks_libur_per_bulan", 3)
	minAvailable := s.configService.GetInt(configs, "min_available_per_hari", 2)

	karyawan, err := s.karyawanRepo.GetByID(karyawanID)
	if err != nil {
		return 0, false, nil, err
	}

	// 1. Check window
	if !s.isWithinWindow(tanggal) {
		return 0, false, nil, ErrOutOfWindow
	}

	// 2. Check quota
	count, err := s.repo.CountEmployeeLeavesInMonth(karyawanID, tanggal.Year(), tanggal.Month())
	if err != nil {
		return 0, false, nil, err
	}
	if int(count) >= maxLibur {
		return 0, false, nil, ErrKuotaHabis
	}

	// 3. Calculate availability
	availableAfter, err := s.calculateAvailabilityAfter(karyawan.TokoID, tanggal, karyawanID)
	if err != nil {
		return 0, false, nil, err
	}

	needsBackup := availableAfter < minAvailable && availableAfter > 0
	var suggestedBackup []model.Karyawan
	if availableAfter < minAvailable {
		suggestedBackup, _ = s.repo.GetAvailableBackups(tanggal, karyawanID)
	}

	return availableAfter, needsBackup, suggestedBackup, nil
}

func (s *JadwalLiburService) CreatePlanned(karyawanID uuid.UUID, tanggal time.Time, backupKaryawanID *uuid.UUID) (*model.JadwalLibur, error) {
	configs, _ := s.configService.GetConfigMap()
	maxLibur := s.configService.GetInt(configs, "maks_libur_per_bulan", 3)
	minAvailable := s.configService.GetInt(configs, "min_available_per_hari", 2)

	karyawan, err := s.karyawanRepo.GetByID(karyawanID)
	if err != nil {
		return nil, err
	}

	if !s.isWithinWindow(tanggal) {
		return nil, ErrOutOfWindow
	}

	count, err := s.repo.CountEmployeeLeavesInMonth(karyawanID, tanggal.Year(), tanggal.Month())
	if err != nil {
		return nil, err
	}
	if int(count) >= maxLibur {
		return nil, ErrKuotaHabis
	}

	availableAfter, err := s.calculateAvailabilityAfter(karyawan.TokoID, tanggal, karyawanID)
	if err != nil {
		return nil, err
	}

	if availableAfter == 0 {
		return nil, ErrNoBackupAvailable
	}

	jadwal := &model.JadwalLibur{
		KaryawanID: karyawanID,
		Tanggal:    tanggal,
		Tipe:       "planned",
	}

	var backup *model.BackupAssignment
	if availableAfter < minAvailable {
		if backupKaryawanID == nil {
			return nil, ErrBackupRequired
		}
		if *backupKaryawanID == karyawanID {
			return nil, ErrBackupSelf
		}
		backupOnLeave, _ := s.isKaryawanOnLeave(*backupKaryawanID, tanggal)
		if backupOnLeave {
			return nil, ErrBackupInvalid
		}

		backup = &model.BackupAssignment{
			BackupKaryawanID: *backupKaryawanID,
			AssignedBy:       karyawanID,
		}
	}

	err = s.repo.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(jadwal).Error; err != nil {
			return err
		}
		if backup != nil {
			backup.JadwalLiburID = jadwal.ID
			if err := tx.Create(backup).Error; err != nil {
				return err
			}
		}

		payload := map[string]interface{}{
			"jadwal": jadwal,
		}
		if backup != nil {
			payload["backup"] = backup
		}

		return s.auditService.Log(tx, &karyawanID, "CREATE_JADWAL_LIBUR", "jadwal_libur", jadwal.ID, payload)
	})

	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(jadwal.ID)
}

func (s *JadwalLiburService) CreateUnplanned(karyawanID uuid.UUID, tanggal time.Time, adminID uuid.UUID) (*model.JadwalLibur, int, []model.Karyawan, error) {
	configs, _ := s.configService.GetConfigMap()
	minAvailable := s.configService.GetInt(configs, "min_available_per_hari", 2)

	karyawan, err := s.karyawanRepo.GetByID(karyawanID)
	if err != nil {
		return nil, 0, nil, err
	}

	jadwal := &model.JadwalLibur{
		KaryawanID: karyawanID,
		Tanggal:    tanggal,
		Tipe:       "unplanned",
	}

	err = s.repo.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(jadwal).Error; err != nil {
			return err
		}
		
		payload := map[string]interface{}{
			"jadwal": jadwal,
		}
		return s.auditService.Log(tx, &adminID, "CREATE_UNPLANNED_LEAVE", "jadwal_libur", jadwal.ID, payload)
	})

	if err != nil {
		return nil, 0, nil, err
	}

	availableAfter, _ := s.calculateAvailabilityAfter(karyawan.TokoID, tanggal, uuid.Nil)

	var suggestedBackup []model.Karyawan
	if availableAfter < minAvailable {
		suggestedBackup, _ = s.repo.GetAvailableBackups(tanggal, karyawanID)
	}

	created, _ := s.repo.GetByID(jadwal.ID)
	return created, availableAfter, suggestedBackup, nil
}

func (s *JadwalLiburService) Update(id uuid.UUID, tanggal time.Time, backupKaryawanID *uuid.UUID, requesterID uuid.UUID) (*model.JadwalLibur, error) {
	configs, _ := s.configService.GetConfigMap()
	maxLibur := s.configService.GetInt(configs, "maks_libur_per_bulan", 3)
	minAvailable := s.configService.GetInt(configs, "min_available_per_hari", 2)

	oldJadwal, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	now := timeutil.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, timeutil.Loc)
	
	if oldJadwal.Tanggal.Before(today) || oldJadwal.Tanggal.Equal(today) {
		return nil, ErrTanggalTerlewat
	}

	if oldJadwal.Tipe == "planned" {
		if !s.isWithinWindow(tanggal) {
			return nil, ErrOutOfWindow
		}
		count, _ := s.repo.CountEmployeeLeavesInMonth(oldJadwal.KaryawanID, tanggal.Year(), tanggal.Month())
		if oldJadwal.Tanggal.Month() != tanggal.Month() || oldJadwal.Tanggal.Year() != tanggal.Year() {
			if int(count) >= maxLibur {
				return nil, ErrKuotaHabis
			}
		}
	}

	availableAfter, _ := s.calculateAvailabilityAfter(oldJadwal.Karyawan.TokoID, tanggal, oldJadwal.KaryawanID)
	if availableAfter == 0 && oldJadwal.Tipe == "planned" {
		return nil, ErrNoBackupAvailable
	}

	data := map[string]interface{}{
		"tanggal": tanggal,
	}

	var backup *model.BackupAssignment
	if availableAfter < minAvailable && oldJadwal.Tipe == "planned" {
		if backupKaryawanID == nil {
			return nil, ErrBackupRequired
		}
		if *backupKaryawanID == oldJadwal.KaryawanID {
			return nil, ErrBackupSelf
		}
		backupOnLeave, _ := s.isKaryawanOnLeave(*backupKaryawanID, tanggal)
		if backupOnLeave {
			return nil, ErrBackupInvalid
		}
		backup = &model.BackupAssignment{
			BackupKaryawanID: *backupKaryawanID,
			AssignedBy:       oldJadwal.KaryawanID,
		}
	}

	err = s.repo.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("jadwal_libur_id = ?", id).Delete(&model.BackupAssignment{}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.JadwalLibur{}).Where("id = ?", id).Updates(data).Error; err != nil {
			return err
		}
		if backup != nil {
			backup.JadwalLiburID = id
			if err := tx.Create(backup).Error; err != nil {
				return err
			}
		}

		payload := map[string]interface{}{
			"before": map[string]interface{}{
				"tanggal": oldJadwal.Tanggal,
				"has_backup": oldJadwal.BackupAssignment != nil,
			},
			"after": map[string]interface{}{
				"tanggal": tanggal,
				"has_backup": backup != nil,
			},
		}

		return s.auditService.Log(tx, &requesterID, "UPDATE_JADWAL_LIBUR", "jadwal_libur", id, payload)
	})

	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

func (s *JadwalLiburService) Delete(id uuid.UUID, requesterID uuid.UUID) error {
	jadwal, err := s.repo.GetByID(id)
	if err != nil {
		return ErrNotFound
	}

	var requester model.Karyawan
	if err := s.repo.DB().First(&requester, "id = ?", requesterID).Error; err != nil {
		return err
	}

	if jadwal.KaryawanID != requesterID && requester.Role != "admin" {
		return ErrUnauthorized
	}

	isAdmin := requester.Role == "admin"
	skipDateValidation := isAdmin && jadwal.Tipe == "unplanned"

	now := timeutil.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, timeutil.Loc)

	if !skipDateValidation {
		if jadwal.Tanggal.Before(today) || jadwal.Tanggal.Equal(today) {
			return ErrTanggalTerlewat
		}
	}

	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		payload := map[string]interface{}{
			"jadwal": jadwal,
		}
		
		action := "DELETE_JADWAL_LIBUR"
		if jadwal.Tipe == "unplanned" {
			action = "DELETE_UNPLANNED_LEAVE"
		}
		
		if err := s.auditService.Log(tx, &requesterID, action, "jadwal_libur", id, payload); err != nil {
			return err
		}

		return tx.Delete(&model.JadwalLibur{}, "id = ?", id).Error
	})
}

// Helpers
func (s *JadwalLiburService) isWithinWindow(tanggal time.Time) bool {
	currentMonth := timeutil.StartOfCurrentMonth()
	nextMonth := currentMonth.AddDate(0, 1, 0)
	lastDayOfNextMonth := nextMonth.AddDate(0, 1, -1)

	// Normalize target date to local midnight
	target := time.Date(tanggal.Year(), tanggal.Month(), tanggal.Day(), 0, 0, 0, 0, timeutil.Loc)
	
	return (target.After(currentMonth) || target.Equal(currentMonth)) && (target.Before(lastDayOfNextMonth) || target.Equal(lastDayOfNextMonth))
}

func (s *JadwalLiburService) calculateAvailabilityAfter(tokoID uuid.UUID, tanggal time.Time, requesterID uuid.UUID) (int, error) {
	total, err := s.repo.GetTotalEmployeesInStore(tokoID)
	if err != nil {
		return 0, err
	}

	onLeaveIDs, err := s.repo.GetEmployeesOnLeave(tokoID, tanggal)
	if err != nil {
		return 0, err
	}

	// Count how many are on leave, excluding the requester if they are already in the list
	leaveCount := 0
	foundRequester := false
	for _, id := range onLeaveIDs {
		if id == requesterID {
			foundRequester = true
		}
		leaveCount++
	}

	if !foundRequester && requesterID != uuid.Nil {
		leaveCount++
	}

	return int(total) - leaveCount, nil
}

func (s *JadwalLiburService) isKaryawanOnLeave(karyawanID uuid.UUID, tanggal time.Time) (bool, error) {
	// Simple check if this specific employee is on leave on that date
	var count int64
	err := s.karyawanRepo.DB().Model(&model.JadwalLibur{}).Where("karyawan_id = ? AND tanggal = ?", karyawanID, tanggal).Count(&count).Error
	return count > 0, err
}
