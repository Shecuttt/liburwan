package handler

import (
	"backend-liburwan/internal/model"
	"backend-liburwan/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TokoHandler struct {
	service *service.TokoService
}

func NewTokoHandler(service *service.TokoService) *TokoHandler {
	return &TokoHandler{service: service}
}

func (h *TokoHandler) GetAll(c *gin.Context) {
	tokos, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, tokos)
}

func (h *TokoHandler) GetByID(c *gin.Context) {
	idStr := c.Param("toko_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "BAD_REQUEST",
			"message": "ID tidak valid",
		})
		return
	}

	toko, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "Toko tidak ditemukan",
		})
		return
	}
	c.JSON(http.StatusOK, toko)
}

func (h *TokoHandler) Create(c *gin.Context) {
	var req struct {
		Nama    string `json:"nama" binding:"required"`
		IsPusat bool   `json:"is_pusat"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	toko := &model.Toko{
		Nama:    req.Nama,
		IsPusat: req.IsPusat,
	}

	if err := h.service.Create(toko); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, toko)
}
