package handler

import (
	"backend-liburwan/internal/lib/timeutil"
	"backend-liburwan/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JadwalLiburHandler struct {
	service *service.JadwalLiburService
}

func NewJadwalLiburHandler(service *service.JadwalLiburService) *JadwalLiburHandler {
	return &JadwalLiburHandler{service: service}
}

func (h *JadwalLiburHandler) GetAll(c *gin.Context) {
	karyawanID := c.Query("karyawan_id")
	tokoID := c.Query("toko_id")
	bulan := c.Query("bulan")

	jadwals, err := h.service.GetAll(karyawanID, tokoID, bulan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_SERVER_ERROR", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jadwals)
}

func (h *JadwalLiburHandler) CheckAvailability(c *gin.Context) {
	tanggalStr := c.Query("tanggal")
	tanggal, err := time.ParseInLocation("2006-01-02", tanggalStr, timeutil.Loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "Format tanggal salah (YYYY-MM-DD)"})
		return
	}

	karyawanIDStr, _ := c.Get("karyawan_id")
	karyawanID, err := uuid.Parse(karyawanIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED", "message": "Invalid context karyawan_id"})
		return
	}

	availableAfter, needsBackup, suggested, err := h.service.CheckAvailability(karyawanID, tanggal)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"available_count_after": availableAfter,
		"needs_backup":         needsBackup,
		"suggested_backup":     suggested,
	})
}

func (h *JadwalLiburHandler) CreatePlanned(c *gin.Context) {
	var req struct {
		Tanggal          string  `json:"tanggal" binding:"required"`
		BackupKaryawanID *string `json:"backup_karyawan_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": err.Error()})
		return
	}

	tanggal, _ := time.ParseInLocation("2006-01-02", req.Tanggal, timeutil.Loc)
	karyawanIDStr, _ := c.Get("karyawan_id")
	karyawanID, _ := uuid.Parse(karyawanIDStr.(string))

	var backupID *uuid.UUID
	if req.BackupKaryawanID != nil && *req.BackupKaryawanID != "" {
		id, _ := uuid.Parse(*req.BackupKaryawanID)
		backupID = &id
	}

	jadwal, err := h.service.CreatePlanned(karyawanID, tanggal, backupID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, jadwal)
}

func (h *JadwalLiburHandler) CreateUnplanned(c *gin.Context) {
	var req struct {
		KaryawanID string `json:"karyawan_id" binding:"required"`
		Tanggal    string `json:"tanggal" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": err.Error()})
		return
	}

	karyawanID, _ := uuid.Parse(req.KaryawanID)
	tanggal, _ := time.ParseInLocation("2006-01-02", req.Tanggal, timeutil.Loc)
	
	adminIDStr, _ := c.Get("karyawan_id")
	adminID, _ := uuid.Parse(adminIDStr.(string))

	jadwal, availableAfter, suggested, err := h.service.CreateUnplanned(karyawanID, tanggal, adminID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"jadwal_libur":       jadwal,
		"availability_after": availableAfter,
		"suggested_backup":   suggested,
	})
}

func (h *JadwalLiburHandler) GetByID(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	jadwal, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "Jadwal libur tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, jadwal)
}

func (h *JadwalLiburHandler) Update(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	var req struct {
		Tanggal          string  `json:"tanggal" binding:"required"`
		BackupKaryawanID *string `json:"backup_karyawan_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": err.Error()})
		return
	}

	tanggal, _ := time.ParseInLocation("2006-01-02", req.Tanggal, timeutil.Loc)
	var backupID *uuid.UUID
	if req.BackupKaryawanID != nil && *req.BackupKaryawanID != "" {
		bid, _ := uuid.Parse(*req.BackupKaryawanID)
		backupID = &bid
	}

	requesterIDStr, _ := c.Get("karyawan_id")
	requesterID, _ := uuid.Parse(requesterIDStr.(string))

	jadwal, err := h.service.Update(id, tanggal, backupID, requesterID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, jadwal)
}

func (h *JadwalLiburHandler) Delete(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	
	requesterIDStr, _ := c.Get("karyawan_id")
	requesterID, _ := uuid.Parse(requesterIDStr.(string))

	if err := h.service.Delete(id, requesterID); err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *JadwalLiburHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrOutOfWindow:
		c.JSON(http.StatusBadRequest, gin.H{"code": "OUT_OF_WINDOW", "message": "Hanya bisa booking bulan berjalan dan bulan depan"})
	case service.ErrKuotaHabis:
		c.JSON(http.StatusBadRequest, gin.H{"code": "KUOTA_HABIS", "message": "Kuota libur Anda untuk bulan ini sudah habis"})
	case service.ErrBackupRequired:
		c.JSON(http.StatusBadRequest, gin.H{"code": "BACKUP_REQUIRED", "message": "Toko akan tersisa 1 orang, wajib sertakan backup_karyawan_id"})
	case service.ErrBackupInvalid:
		c.JSON(http.StatusBadRequest, gin.H{"code": "BACKUP_INVALID", "message": "Karyawan yang dipilih sebagai backup juga libur di tanggal ini"})
	case service.ErrBackupSelf:
		c.JSON(http.StatusBadRequest, gin.H{"code": "BACKUP_INVALID", "message": "Tidak dapat menunjuk diri sendiri sebagai backup"})
	case service.ErrNoBackupAvailable:
		c.JSON(http.StatusBadRequest, gin.H{"code": "NO_BACKUP_AVAILABLE", "message": "Tidak ada karyawan lain yang bisa jadi backup di tanggal ini"})
	case service.ErrTanggalTerlewat:
		c.JSON(http.StatusBadRequest, gin.H{"code": "TANGGAL_TERLEWAT", "message": "Tidak bisa memproses jadwal libur yang tanggalnya sudah lewat"})
	case service.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "Jadwal libur tidak ditemukan"})
	case service.ErrUnauthorized:
		c.JSON(http.StatusForbidden, gin.H{"code": "UNAUTHORIZED_ACTION", "message": "Anda tidak memiliki akses untuk melakukan tindakan ini"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_SERVER_ERROR", "message": err.Error()})
	}
}
