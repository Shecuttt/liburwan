package handler

import (
	"backend-liburwan/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BackupAssignmentHandler struct {
	service *service.BackupAssignmentService
}

func NewBackupAssignmentHandler(service *service.BackupAssignmentService) *BackupAssignmentHandler {
	return &BackupAssignmentHandler{service: service}
}

func (h *BackupAssignmentHandler) Create(c *gin.Context) {
	var req struct {
		JadwalLiburID    string `json:"jadwal_libur_id" binding:"required"`
		BackupKaryawanID string `json:"backup_karyawan_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": err.Error()})
		return
	}

	jadwalLiburID, _ := uuid.Parse(req.JadwalLiburID)
	backupKaryawanID, _ := uuid.Parse(req.BackupKaryawanID)
	
	assignedByStr, _ := c.Get("karyawan_id")
	assignedBy, _ := uuid.Parse(assignedByStr.(string))

	backup, err := h.service.Create(jadwalLiburID, backupKaryawanID, assignedBy)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, backup)
}

func (h *BackupAssignmentHandler) Delete(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))

	deletedByStr, _ := c.Get("karyawan_id")
	deletedBy, _ := uuid.Parse(deletedByStr.(string))

	if err := h.service.Delete(id, deletedBy); err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *BackupAssignmentHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "Record tidak ditemukan"})
	case service.ErrAlreadyAssigned:
		c.JSON(http.StatusBadRequest, gin.H{"code": "ALREADY_ASSIGNED", "message": "Jadwal libur ini sudah memiliki backup assignment"})
	case service.ErrBackupInvalid:
		c.JSON(http.StatusBadRequest, gin.H{"code": "BACKUP_INVALID", "message": "Karyawan yang dipilih sebagai backup juga libur di tanggal ini"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_SERVER_ERROR", "message": err.Error()})
	}
}
