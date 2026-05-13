package service

import (
	"backend-liburwan/internal/model"
	"backend-liburwan/internal/repository"
	"time"

	"github.com/google/uuid"
)

type KalenderService struct {
	jadwalRepo   *repository.JadwalLiburRepository
	tokoRepo     *repository.TokoRepository
	karyawanRepo *repository.KaryawanRepository
}

func NewKalenderService(jadwalRepo *repository.JadwalLiburRepository, tokoRepo *repository.TokoRepository, karyawanRepo *repository.KaryawanRepository) *KalenderService {
	return &KalenderService{
		jadwalRepo:   jadwalRepo,
		tokoRepo:     tokoRepo,
		karyawanRepo: karyawanRepo,
	}
}

type CalendarResponse struct {
	Bulan string        `json:"bulan"`
	Hari  []CalendarDay `json:"hari"`
}

type CalendarDay struct {
	Tanggal string         `json:"tanggal"`
	Toko    []CalendarToko `json:"toko"`
}

type CalendarToko struct {
	TokoID         uuid.UUID       `json:"toko_id"`
	TokoNama       string          `json:"toko_nama"`
	AvailableCount int             `json:"available_count"`
	IsRawan        bool            `json:"is_rawan"`
	Libur          []CalendarLibur `json:"libur"`
}

type CalendarLibur struct {
	ID                 uuid.UUID  `json:"id"`
	KaryawanID         uuid.UUID  `json:"karyawan_id"`
	KaryawanNama       string     `json:"karyawan_nama"`
	Tipe               string     `json:"tipe"`
	BackupAssignmentID *uuid.UUID `json:"backup_assignment_id"`
	BackupKaryawanNama *string    `json:"backup_karyawan_nama"`
}

func (s *KalenderService) GetCalendar(bulan string) (*CalendarResponse, error) {
	// 1. Parse bulan (YYYY-MM)
	t, err := time.ParseInLocation("2006-01", bulan, time.Local)
	if err != nil {
		return nil, err
	}

	// 2. Window validation (current and next month)
	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	nextMonth := currentMonth.AddDate(0, 1, 0)
	
	// Compare only year and month
	targetMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	
	if !targetMonth.Equal(currentMonth) && !targetMonth.Equal(nextMonth) {
		return nil, ErrOutOfWindow
	}

	// 3. Get all toko
	tokos, err := s.tokoRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// 4. Get all karyawan grouped by toko
	karyawans, err := s.karyawanRepo.GetAll("")
	if err != nil {
		return nil, err
	}
	tokoKaryawanMap := make(map[uuid.UUID][]model.Karyawan)
	for _, k := range karyawans {
		tokoKaryawanMap[k.TokoID] = append(tokoKaryawanMap[k.TokoID], k)
	}

	// 5. Get all jadwal libur for the month
	jadwals, err := s.jadwalRepo.GetAll("", "", bulan)
	if err != nil {
		return nil, err
	}

	// Group jadwals by date and toko
	dateTokoLiburMap := make(map[string]map[uuid.UUID][]model.JadwalLibur)
	for _, j := range jadwals {
		dateStr := j.Tanggal.Format("2006-01-02")
		if dateTokoLiburMap[dateStr] == nil {
			dateTokoLiburMap[dateStr] = make(map[uuid.UUID][]model.JadwalLibur)
		}
		dateTokoLiburMap[dateStr][j.Karyawan.TokoID] = append(dateTokoLiburMap[dateStr][j.Karyawan.TokoID], j)
	}

	// 6. Generate days for the month
	daysInMonth := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()
	resp := &CalendarResponse{
		Bulan: bulan,
		Hari:  make([]CalendarDay, 0, daysInMonth),
	}

	for d := 1; d <= daysInMonth; d++ {
		currDate := time.Date(t.Year(), t.Month(), d, 0, 0, 0, 0, time.Local)
		dateStr := currDate.Format("2006-01-02")
		
		calDay := CalendarDay{
			Tanggal: dateStr,
			Toko:    make([]CalendarToko, 0, len(tokos)),
		}

		for _, toko := range tokos {
			liburs := dateTokoLiburMap[dateStr][toko.ID]
			totalKaryawan := len(tokoKaryawanMap[toko.ID])
			availableCount := totalKaryawan - len(liburs)

			calLiburs := make([]CalendarLibur, 0, len(liburs))
			for _, l := range liburs {
				var backupNama *string
				var backupID *uuid.UUID
				if l.BackupAssignment != nil {
					backupNama = &l.BackupAssignment.BackupKaryawan.Nama
					backupID = &l.BackupAssignment.ID
				}
				calLiburs = append(calLiburs, CalendarLibur{
					ID:                 l.ID,
					KaryawanID:         l.KaryawanID,
					KaryawanNama:       l.Karyawan.Nama,
					Tipe:               l.Tipe,
					BackupAssignmentID: backupID,
					BackupKaryawanNama: backupNama,
				})
			}

			calDay.Toko = append(calDay.Toko, CalendarToko{
				TokoID:         toko.ID,
				TokoNama:       toko.Nama,
				AvailableCount: availableCount,
				IsRawan:        availableCount <= 2,
				Libur:          calLiburs,
			})
		}
		resp.Hari = append(resp.Hari, calDay)
	}

	return resp, nil
}
