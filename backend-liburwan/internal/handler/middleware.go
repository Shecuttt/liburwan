package handler

import (
	"backend-liburwan/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED", "message": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED", "message": "Invalid authorization format"})
			c.Abort()
			return
		}

		claims, err := authService.ValidateJWT(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED", "message": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("karyawan_id", claims["karyawan_id"])
		c.Set("role", claims["role"])
		c.Set("toko_id", claims["toko_id"])

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    "FORBIDDEN",
				"message": "Akses ditolak, hanya untuk admin",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
