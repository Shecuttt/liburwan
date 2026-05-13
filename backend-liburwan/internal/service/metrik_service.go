package service

import (
	"backend-liburwan/internal/repository"
	"strconv"

	"github.com/google/uuid"
)

type MetrikService struct {
	repo         *repository.MetrikRepository
	karyawanRepo *repository.KaryawanRepository
	tokoRepo     *repository.TokoRepository
	jadwalRepo   *repository.JadwalLiburRepository
}

func NewMetrikService(
	repo *repository.MetrikRepository,
	karyawanRepo *repository.KaryawanRepository,
	tokoRepo *repository.TokoRepository,
	jadwalRepo *repository.JadwalLiburRepository,
) *MetrikService {
	return &MetrikService{
		repo:         repo,
		karyawanRepo: karyawanRepo,
		tokoRepo:     tokoRepo,
		jadwalRepo:   jadwalRepo,
	}
}

func (s *MetrikService) GetKaryawanMetrik(karyawanID uuid.UUID, bulan string) (map[string]interface{}, error) {
	karyawan, err := s.karyawanRepo.GetByID(karyawanID)
	if err != nil {
		return nil, err
	}

	totalPlanned, _ := s.repo.CountLeaves(karyawanID, "planned", bulan)
	totalUnplanned, _ := s.repo.CountLeaves(karyawanID, "unplanned", bulan)
	totalJadiBackup, _ := s.repo.CountAsBackup(karyawanID, bulan)

	histori, _ := s.jadwalRepo.GetAll(karyawanID.String(), "", bulan)

	return map[string]interface{}{
		"karyawan": karyawan,
		"ringkasan": map[string]interface{}{
			"total_planned":   totalPlanned,
			"total_unplanned": totalUnplanned,
			"total_jadi_backup": totalJadiBackup,
		},
		"histori": histori,
	}, nil
}

func (s *MetrikService) GetTokoMetrik(tokoID uuid.UUID, bulan string) (map[string]interface{}, error) {
	toko, err := s.tokoRepo.GetByID(tokoID)
	if err != nil {
		return nil, err
	}

	totalKaryawan, _ := s.karyawanRepo.CountInStore(tokoID)
	minAvailableStr, _ := s.repo.GetConfigValue("min_available_per_hari")
	minAvailable, _ := strconv.Atoi(minAvailableStr)

	totalUnplanned, _ := s.repo.CountStoreUnplanned(tokoID, bulan)
	totalBackupDariLuar, _ := s.repo.CountBackupFromOutside(tokoID, bulan)

	onLeavePerDay, _ := s.repo.GetAvailabilityPerDay(tokoID, bulan)

	var hariKritis []map[string]interface{}
	for date, leaveCount := range onLeavePerDay {
		available := int(totalKaryawan) - leaveCount
		if available <= minAvailable {
			hariKritis = append(hariKritis, map[string]interface{}{
				"tanggal":         date,
				"available_count": available,
			})
		}
	}

	return map[string]interface{}{
		"toko": toko,
		"ringkasan": map[string]interface{}{
			"total_hari_kritis":       len(hariKritis),
			"total_backup_dari_luar": totalBackupDariLuar,
			"total_unplanned":        totalUnplanned,
		},
		"hari_kritis": hariKritis,
	}, nil
}
