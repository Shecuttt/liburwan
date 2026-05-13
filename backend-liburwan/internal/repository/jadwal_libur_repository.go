package repository

import (
	"backend-liburwan/internal/lib/timeutil"
	"backend-liburwan/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JadwalLiburRepository struct {
	db *gorm.DB
}

func NewJadwalLiburRepository(db *gorm.DB) *JadwalLiburRepository {
	return &JadwalLiburRepository{db: db}
}

func (r *JadwalLiburRepository) DB() *gorm.DB {
	return r.db
}

func (r *JadwalLiburRepository) GetAll(karyawanID, tokoID, bulan string) ([]model.JadwalLibur, error) {
	var jadwals []model.JadwalLibur
	query := r.db.Preload("Karyawan.Toko").Preload("BackupAssignment.BackupKaryawan").Preload("BackupAssignment.Assigner")

	if karyawanID != "" {
		query = query.Where("karyawan_id = ?", karyawanID)
	}
	if tokoID != "" {
		query = query.Joins("JOIN karyawan ON karyawan.id = jadwal_libur.karyawan_id").Where("karyawan.toko_id = ?", tokoID)
	}
	if bulan != "" {
		// Expecting YYYY-MM
		t, err := time.ParseInLocation("2006-01", bulan, timeutil.Loc)
		if err == nil {
			start := t
			end := t.AddDate(0, 1, 0)
			query = query.Where("tanggal >= ? AND tanggal < ?", start, end)
		}
	}

	err := query.Find(&jadwals).Error
	return jadwals, err
}

func (r *JadwalLiburRepository) GetByID(id uuid.UUID) (*model.JadwalLibur, error) {
	var jadwal model.JadwalLibur
	err := r.db.Preload("Karyawan.Toko").Preload("BackupAssignment.BackupKaryawan").Preload("BackupAssignment.Assigner").First(&jadwal, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &jadwal, nil
}

func (r *JadwalLiburRepository) Create(jadwal *model.JadwalLibur) error {
	return r.db.Create(jadwal).Error
}

func (r *JadwalLiburRepository) CreateWithBackup(jadwal *model.JadwalLibur, backup *model.BackupAssignment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(jadwal).Error; err != nil {
			return err
		}
		backup.JadwalLiburID = jadwal.ID
		if err := tx.Create(backup).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *JadwalLiburRepository) Update(id uuid.UUID, data map[string]interface{}) error {
	return r.db.Model(&model.JadwalLibur{}).Where("id = ?", id).Updates(data).Error
}

func (r *JadwalLiburRepository) UpdateWithBackup(id uuid.UUID, data map[string]interface{}, backup *model.BackupAssignment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing backup if any
		if err := tx.Where("jadwal_libur_id = ?", id).Delete(&model.BackupAssignment{}).Error; err != nil {
			return err
		}
		// Update jadwal
		if err := tx.Model(&model.JadwalLibur{}).Where("id = ?", id).Updates(data).Error; err != nil {
			return err
		}
		// Create new backup if provided
		if backup != nil {
			backup.JadwalLiburID = id
			if err := tx.Create(backup).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *JadwalLiburRepository) Delete(id uuid.UUID) error {
	// Cascade delete is handled by DB constraint (ON DELETE CASCADE)
	return r.db.Delete(&model.JadwalLibur{}, "id = ?", id).Error
}

func (r *JadwalLiburRepository) CountEmployeeLeavesInMonth(karyawanID uuid.UUID, year int, month time.Month) (int64, error) {
	var count int64
	start := time.Date(year, month, 1, 0, 0, 0, 0, timeutil.Loc)
	end := start.AddDate(0, 1, 0)
	err := r.db.Model(&model.JadwalLibur{}).
		Where("karyawan_id = ? AND tanggal >= ? AND tanggal < ? AND tipe = 'planned'", karyawanID, start, end).
		Count(&count).Error
	return count, err
}

func (r *JadwalLiburRepository) GetEmployeesOnLeave(tokoID uuid.UUID, tanggal time.Time) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := r.db.Model(&model.JadwalLibur{}).
		Joins("JOIN karyawan ON karyawan.id = jadwal_libur.karyawan_id").
		Where("karyawan.toko_id = ? AND jadwal_libur.tanggal = ?", tokoID, tanggal).
		Pluck("karyawan_id", &ids).Error
	return ids, err
}

func (r *JadwalLiburRepository) GetAvailableBackups(tanggal time.Time, excludeKaryawanID uuid.UUID) ([]model.Karyawan, error) {
	var karyawans []model.Karyawan
	// Employees who are NOT on leave and NOT the requester
	subQuery := r.db.Model(&model.JadwalLibur{}).Select("karyawan_id").Where("tanggal = ?", tanggal)
	err := r.db.Preload("Toko").
		Where("id NOT IN (?) AND id <> ?", subQuery, excludeKaryawanID).
		Find(&karyawans).Error
	return karyawans, err
}

func (r *JadwalLiburRepository) GetTotalEmployeesInStore(tokoID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.Karyawan{}).Where("toko_id = ?", tokoID).Count(&count).Error
	return count, err
}
