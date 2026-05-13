package handler

import (
	"backend-liburwan/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MetrikHandler struct {
	service *service.MetrikService
}

func NewMetrikHandler(service *service.MetrikService) *MetrikHandler {
	return &MetrikHandler{service: service}
}

func (h *MetrikHandler) GetKaryawanMetrik(c *gin.Context) {
	karyawanID, err := uuid.Parse(c.Param("karyawan_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "Invalid karyawan_id"})
		return
	}

	bulan := c.Query("bulan")
	data, err := h.service.GetKaryawanMetrik(karyawanID, bulan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_SERVER_ERROR", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *MetrikHandler) GetTokoMetrik(c *gin.Context) {
	tokoID, err := uuid.Parse(c.Param("toko_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "Invalid toko_id"})
		return
	}

	bulan := c.Query("bulan")
	data, err := h.service.GetTokoMetrik(tokoID, bulan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_SERVER_ERROR", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
