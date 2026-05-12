package handler

import (
	"backend-liburwan/internal/model"
	"backend-liburwan/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type KaryawanHandler struct {
	service *service.KaryawanService
}

func NewKaryawanHandler(service *service.KaryawanService) *KaryawanHandler {
	return &KaryawanHandler{service: service}
}

type KaryawanResponse struct {
	ID        uuid.UUID `json:"id"`
	Nama      string    `json:"nama"`
	Role      string    `json:"role"`
	TokoID    uuid.UUID `json:"toko_id"`
	TokoNama  string    `json:"toko_nama"`
	CreatedAt time.Time `json:"created_at"`
}

func mapToKaryawanResponse(k model.Karyawan) KaryawanResponse {
	return KaryawanResponse{
		ID:        k.ID,
		Nama:      k.Nama,
		Role:      k.Role,
		TokoID:    k.TokoID,
		TokoNama:  k.Toko.Nama,
		CreatedAt: k.CreatedAt,
	}
}

func (h *KaryawanHandler) GetAll(c *gin.Context) {
	tokoID := c.Query("toko_id")
	karyawans, err := h.service.GetAll(tokoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}

	res := make([]KaryawanResponse, len(karyawans))
	for i, k := range karyawans {
		res[i] = mapToKaryawanResponse(k)
	}

	c.JSON(http.StatusOK, res)
}

func (h *KaryawanHandler) GetByID(c *gin.Context) {
	idStr := c.Param("karyawan_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "BAD_REQUEST",
			"message": "ID tidak valid",
		})
		return
	}

	karyawan, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "Karyawan tidak ditemukan",
		})
		return
	}
	c.JSON(http.StatusOK, mapToKaryawanResponse(*karyawan))
}

func (h *KaryawanHandler) Create(c *gin.Context) {
	var req struct {
		Nama   string `json:"nama" binding:"required"`
		Role   string `json:"role" binding:"required"`
		TokoID string `json:"toko_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	tokoUUID, err := uuid.Parse(req.TokoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "BAD_REQUEST",
			"message": "Toko ID tidak valid",
		})
		return
	}

	karyawan := &model.Karyawan{
		Nama:   req.Nama,
		Role:   req.Role,
		TokoID: tokoUUID,
	}

	if err := h.service.Create(karyawan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}

	// Fetch created karyawan to get Toko data for response
	createdKaryawan, _ := h.service.GetByID(karyawan.ID)
	c.JSON(http.StatusCreated, mapToKaryawanResponse(*createdKaryawan))
}

func (h *KaryawanHandler) Update(c *gin.Context) {
	idStr := c.Param("karyawan_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "BAD_REQUEST",
			"message": "ID tidak valid",
		})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	if err := h.service.Update(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}

	updatedKaryawan, _ := h.service.GetByID(id)
	c.JSON(http.StatusOK, mapToKaryawanResponse(*updatedKaryawan))
}
