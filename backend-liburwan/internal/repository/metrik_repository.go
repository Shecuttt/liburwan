package repository

import (
	"backend-liburwan/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetrikRepository struct {
	db *gorm.DB
}

func NewMetrikRepository(db *gorm.DB) *MetrikRepository {
	return &MetrikRepository{db: db}
}

func (r *MetrikRepository) CountLeaves(karyawanID uuid.UUID, tipe string, bulan string) (int64, error) {
	var count int64
	query := r.db.Model(&model.JadwalLibur{}).Where("karyawan_id = ? AND tipe = ?", karyawanID, tipe)
	if bulan != "" {
		start := bulan + "-01"
		query = query.Where("tanggal >= ? AND tanggal < ?::date + interval '1 month'", start, start)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *MetrikRepository) CountAsBackup(karyawanID uuid.UUID, bulan string) (int64, error) {
	var count int64
	query := r.db.Model(&model.BackupAssignment{}).
		Joins("JOIN jadwal_libur ON jadwal_libur.id = backup_assignment.jadwal_libur_id").
		Where("backup_karyawan_id = ?", karyawanID)
	if bulan != "" {
		start := bulan + "-01"
		query = query.Where("jadwal_libur.tanggal >= ? AND jadwal_libur.tanggal < ?::date + interval '1 month'", start, start)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *MetrikRepository) CountStoreUnplanned(tokoID uuid.UUID, bulan string) (int64, error) {
	var count int64
	query := r.db.Model(&model.JadwalLibur{}).
		Joins("JOIN karyawan ON karyawan.id = jadwal_libur.karyawan_id").
		Where("karyawan.toko_id = ? AND tipe = 'unplanned'", tokoID)
	if bulan != "" {
		start := bulan + "-01"
		query = query.Where("tanggal >= ? AND tanggal < ?::date + interval '1 month'", start, start)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *MetrikRepository) CountBackupFromOutside(tokoID uuid.UUID, bulan string) (int64, error) {
	var count int64
	query := r.db.Model(&model.BackupAssignment{}).
		Joins("JOIN jadwal_libur ON jadwal_libur.id = backup_assignment.jadwal_libur_id").
		Joins("JOIN karyawan k_owner ON k_owner.id = jadwal_libur.karyawan_id").
		Joins("JOIN karyawan k_backup ON k_backup.id = backup_assignment.backup_karyawan_id").
		Where("k_owner.toko_id = ? AND k_backup.toko_id <> ?", tokoID, tokoID)
	if bulan != "" {
		start := bulan + "-01"
		query = query.Where("jadwal_libur.tanggal >= ? AND jadwal_libur.tanggal < ?::date + interval '1 month'", start, start)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *MetrikRepository) GetConfigValue(key string) (string, error) {
	var config model.Konfigurasi
	err := r.db.Where("key = ?", key).First(&config).Error
	return config.Value, err
}

func (r *MetrikRepository) GetAvailabilityPerDay(tokoID uuid.UUID, bulan string) (map[string]int, error) {
	// Returns a map of date string -> count of employees ON LEAVE
	// This will be subtracted from total employees in service layer
	type Result struct {
		Tanggal string
		Count   int
	}
	var results []Result
	query := r.db.Model(&model.JadwalLibur{}).
		Select("tanggal::text as tanggal, count(*) as count").
		Joins("JOIN karyawan ON karyawan.id = jadwal_libur.karyawan_id").
		Where("karyawan.toko_id = ?", tokoID).
		Group("tanggal")
	
	if bulan != "" {
		start := bulan + "-01"
		query = query.Where("tanggal >= ? AND tanggal < ?::date + interval '1 month'", start, start)
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	m := make(map[string]int)
	for _, res := range results {
		// PostgreSQL might return YYYY-MM-DD
		m[res.Tanggal] = res.Count
	}
	return m, nil
}
