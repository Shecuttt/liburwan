package service

import (
	"backend-liburwan/internal/model"
	"backend-liburwan/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrOutOfWindow       = errors.New("OUT_OF_WINDOW")
	ErrKuotaHabis        = errors.New("KUOTA_HABIS")
	ErrBackupRequired    = errors.New("BACKUP_REQUIRED")
	ErrBackupInvalid     = errors.New("BACKUP_INVALID")
	ErrNoBackupAvailable = errors.New("NO_BACKUP_AVAILABLE")
	ErrTanggalTerlewat   = errors.New("TANGGAL_TERLEWAT")
	ErrNotFound          = errors.New("NOT_FOUND")
)

type JadwalLiburService struct {
	repo         *repository.JadwalLiburRepository
	karyawanRepo *repository.KaryawanRepository
}

func NewJadwalLiburService(repo *repository.JadwalLiburRepository, karyawanRepo *repository.KaryawanRepository) *JadwalLiburService {
	return &JadwalLiburService{repo: repo, karyawanRepo: karyawanRepo}
}

func (s *JadwalLiburService) GetAll(karyawanID, tokoID, bulan string) ([]model.JadwalLibur, error) {
	return s.repo.GetAll(karyawanID, tokoID, bulan)
}

func (s *JadwalLiburService) GetByID(id uuid.UUID) (*model.JadwalLibur, error) {
	return s.repo.GetByID(id)
}

func (s *JadwalLiburService) CheckAvailability(karyawanID uuid.UUID, tanggal time.Time) (int, bool, []model.Karyawan, error) {
	karyawan, err := s.karyawanRepo.GetByID(karyawanID)
	if err != nil {
		return 0, false, nil, err
	}

	// 1. Check window (Hanya untuk planned, tapi check helper biasanya dipanggil untuk planned)
	if !s.isWithinWindow(tanggal) {
		return 0, false, nil, ErrOutOfWindow
	}

	// 2. Check quota
	count, err := s.repo.CountEmployeeLeavesInMonth(karyawanID, tanggal.Year(), tanggal.Month())
	if err != nil {
		return 0, false, nil, err
	}
	if count >= 3 {
		return 0, false, nil, ErrKuotaHabis
	}

	// 3. Calculate availability
	availableAfter, err := s.calculateAvailabilityAfter(karyawan.TokoID, tanggal, karyawanID)
	if err != nil {
		return 0, false, nil, err
	}

	needsBackup := availableAfter == 1
	var suggestedBackup []model.Karyawan
	if availableAfter < 2 {
		suggestedBackup, _ = s.repo.GetAvailableBackups(tanggal, karyawanID)
	}

	return availableAfter, needsBackup, suggestedBackup, nil
}

func (s *JadwalLiburService) CreatePlanned(karyawanID uuid.UUID, tanggal time.Time, backupKaryawanID *uuid.UUID) (*model.JadwalLibur, error) {
	karyawan, err := s.karyawanRepo.GetByID(karyawanID)
	if err != nil {
		return nil, err
	}

	// 1. Check window
	if !s.isWithinWindow(tanggal) {
		return nil, ErrOutOfWindow
	}

	// 2. Check quota
	count, err := s.repo.CountEmployeeLeavesInMonth(karyawanID, tanggal.Year(), tanggal.Month())
	if err != nil {
		return nil, err
	}
	if count >= 3 {
		return nil, ErrKuotaHabis
	}

	// 3. Calculate availability
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

	if availableAfter == 1 {
		if backupKaryawanID == nil {
			return nil, ErrBackupRequired
		}
		// Check if backup is valid (not on leave)
		backupOnLeave, _ := s.isKaryawanOnLeave( *backupKaryawanID, tanggal)
		if backupOnLeave {
			return nil, ErrBackupInvalid
		}

		backup := &model.BackupAssignment{
			BackupKaryawanID: *backupKaryawanID,
			AssignedBy:       karyawanID, // Self-assigned in planned leave context
		}
		err = s.repo.CreateWithBackup(jadwal, backup)
	} else {
		err = s.repo.Create(jadwal)
	}

	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(jadwal.ID)
}

func (s *JadwalLiburService) CreateUnplanned(karyawanID uuid.UUID, tanggal time.Time) (*model.JadwalLibur, int, []model.Karyawan, error) {
	karyawan, err := s.karyawanRepo.GetByID(karyawanID)
	if err != nil {
		return nil, 0, nil, err
	}

	jadwal := &model.JadwalLibur{
		KaryawanID: karyawanID,
		Tanggal:    tanggal,
		Tipe:       "unplanned",
	}

	if err := s.repo.Create(jadwal); err != nil {
		return nil, 0, nil, err
	}

	// Calculate availability after
	availableAfter, _ := s.calculateAvailabilityAfter(karyawan.TokoID, tanggal, uuid.Nil) // Already created, so just count
	
	var suggestedBackup []model.Karyawan
	if availableAfter < 2 {
		suggestedBackup, _ = s.repo.GetAvailableBackups(tanggal, karyawanID)
	}

	created, _ := s.repo.GetByID(jadwal.ID)
	return created, availableAfter, suggestedBackup, nil
}

func (s *JadwalLiburService) Update(id uuid.UUID, tanggal time.Time, backupKaryawanID *uuid.UUID, requesterRole string) (*model.JadwalLibur, error) {
	oldJadwal, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	// Constraint: old date must not have passed
	if oldJadwal.Tanggal.Before(time.Now().Truncate(24 * time.Hour)) || oldJadwal.Tanggal.Equal(time.Now().Truncate(24 * time.Hour)) {
		return nil, ErrTanggalTerlewat
	}

	// Run all validations for new date (if it's planned)
	if oldJadwal.Tipe == "planned" {
		if !s.isWithinWindow(tanggal) {
			return nil, ErrOutOfWindow
		}
		// Quota: count excluding THIS record if it's in the same month
		// Actually, simpler: if month changes, check quota. If same month, count will include this one, so it should be <= 3.
		count, _ := s.repo.CountEmployeeLeavesInMonth(oldJadwal.KaryawanID, tanggal.Year(), tanggal.Month())
		if oldJadwal.Tanggal.Month() != tanggal.Month() || oldJadwal.Tanggal.Year() != tanggal.Year() {
			if count >= 3 {
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
	if availableAfter == 1 && oldJadwal.Tipe == "planned" {
		if backupKaryawanID == nil {
			return nil, ErrBackupRequired
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

	if err := s.repo.UpdateWithBackup(id, data, backup); err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

func (s *JadwalLiburService) Delete(id uuid.UUID) error {
	jadwal, err := s.repo.GetByID(id)
	if err != nil {
		return ErrNotFound
	}

	if jadwal.Tanggal.Before(time.Now().Truncate(24 * time.Hour)) || jadwal.Tanggal.Equal(time.Now().Truncate(24 * time.Hour)) {
		return ErrTanggalTerlewat
	}

	return s.repo.Delete(id)
}

// Helpers
func (s *JadwalLiburService) isWithinWindow(tanggal time.Time) bool {
	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	nextMonth := currentMonth.AddDate(0, 1, 0)
	lastDayOfNextMonth := nextMonth.AddDate(0, 1, -1)

	// Normalize target date to local midnight
	target := time.Date(tanggal.Year(), tanggal.Month(), tanggal.Day(), 0, 0, 0, 0, time.Local)
	
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
