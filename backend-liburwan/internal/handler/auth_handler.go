package handler

import (
	"backend-liburwan/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"os"
)

type AuthHandler struct {
	authService     *service.AuthService
	karyawanService *service.KaryawanService
}

func NewAuthHandler(authService *service.AuthService, karyawanService *service.KaryawanService) *AuthHandler {
	return &AuthHandler{authService: authService, karyawanService: karyawanService}
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	// In production, use a more secure state (e.g. CSRF token)
	url := h.authService.GetAuthURL("random-state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "Code is required"})
		return
	}

	token, _, err := h.authService.HandleGoogleCallback(code)
	if err != nil {
		if err.Error() == "USER_NOT_REGISTERED" {
			c.JSON(http.StatusForbidden, gin.H{"code": "FORBIDDEN", "message": "Email tidak terdaftar di sistem"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_SERVER_ERROR", "message": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"/auth/callback?token="+token)
}

func (h *AuthHandler) Me(c *gin.Context) {
	karyawanIDVal, exists := c.Get("karyawan_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED", "message": "Karyawan ID tidak ditemukan dalam konteks"})
		return
	}

	karyawanIDStr, ok := karyawanIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_SERVER_ERROR", "message": "Format Karyawan ID tidak valid"})
		return
	}

	karyawanID, err := uuid.Parse(karyawanIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "Format UUID tidak valid"})
		return
	}

	karyawan, err := h.karyawanService.GetByID(karyawanID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "Karyawan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, karyawan)
}
