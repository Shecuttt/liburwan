package handler

import (
	"backend-liburwan/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type KonfigurasiHandler struct {
	service *service.KonfigurasiService
}

func NewKonfigurasiHandler(service *service.KonfigurasiService) *KonfigurasiHandler {
	return &KonfigurasiHandler{service: service}
}

func (h *KonfigurasiHandler) GetAll(c *gin.Context) {
	configs, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_SERVER_ERROR", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, configs)
}

func (h *KonfigurasiHandler) Update(c *gin.Context) {
	key := c.Param("key")
	var req struct {
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": err.Error()})
		return
	}

	adminIDStr, _ := c.Get("karyawan_id")
	adminID, _ := uuid.Parse(adminIDStr.(string))

	if err := h.service.Update(key, req.Value, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_SERVER_ERROR", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Konfigurasi updated successfully"})
}
