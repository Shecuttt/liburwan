package handler

import (
	"backend-liburwan/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type KalenderHandler struct {
	service *service.KalenderService
}

func NewKalenderHandler(service *service.KalenderService) *KalenderHandler {
	return &KalenderHandler{service: service}
}

func (h *KalenderHandler) GetCalendar(c *gin.Context) {
	bulan := c.Query("bulan")
	if bulan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "Parameter bulan (YYYY-MM) wajib disertakan"})
		return
	}

	resp, err := h.service.GetCalendar(bulan)
	if err != nil {
		if err == service.ErrOutOfWindow {
			c.JSON(http.StatusBadRequest, gin.H{"code": "OUT_OF_WINDOW", "message": "Hanya bisa melihat bulan berjalan dan bulan depan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_SERVER_ERROR", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
