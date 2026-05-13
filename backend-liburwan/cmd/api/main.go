package main

import (
	"log"

	"backend-liburwan/internal/config"
	"backend-liburwan/internal/handler"
	"backend-liburwan/internal/lib/timeutil"
	"backend-liburwan/internal/repository"
	"backend-liburwan/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Set global timezone
	time.Local = timeutil.Loc
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Init DB (GORM)
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get *sql.DB from GORM for migrations
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from GORM: %v", err)
	}

	// Run migrations
	if err := config.RunMigrations(sqlDB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Repositories
	tokoRepo := repository.NewTokoRepository(db)
	karyawanRepo := repository.NewKaryawanRepository(db)
	jadwalRepo := repository.NewJadwalLiburRepository(db)
	backupRepo := repository.NewBackupAssignmentRepository(db)
	metrikRepo := repository.NewMetrikRepository(db)
	configRepo := repository.NewKonfigurasiRepository(db)

	// Services
	tokoService := service.NewTokoService(tokoRepo)
	karyawanService := service.NewKaryawanService(karyawanRepo)
	auditLogService := service.NewAuditLogService()
	configService := service.NewKonfigurasiService(configRepo, auditLogService)
	jadwalService := service.NewJadwalLiburService(jadwalRepo, karyawanRepo, configService, auditLogService)
	backupService := service.NewBackupAssignmentService(backupRepo, jadwalRepo, auditLogService)
	metrikService := service.NewMetrikService(metrikRepo, karyawanRepo, tokoRepo, jadwalRepo)
	authService := service.NewAuthService(karyawanRepo)
	kalenderService := service.NewKalenderService(jadwalRepo, tokoRepo, karyawanRepo)

	// Handlers
	tokoHandler := handler.NewTokoHandler(tokoService)
	karyawanHandler := handler.NewKaryawanHandler(karyawanService)
	jadwalHandler := handler.NewJadwalLiburHandler(jadwalService)
	backupHandler := handler.NewBackupAssignmentHandler(backupService)
	metrikHandler := handler.NewMetrikHandler(metrikService)
	configHandler := handler.NewKonfigurasiHandler(configService)
	authHandler := handler.NewAuthHandler(authService, karyawanService)
	kalenderHandler := handler.NewKalenderHandler(kalenderService)

	// Init Gin Router
	r := gin.Default()

	// Apply CORS
	r.Use(handler.CORSMiddleware())

	// API Routes
	api := r.Group("/api/v1")
	{
		// Auth Routes (Public)
		auth := api.Group("/auth")
		{
			auth.GET("/google", authHandler.GoogleLogin)
			auth.GET("/google/callback", authHandler.GoogleCallback)
			auth.GET("/me", handler.JWTMiddleware(authService), authHandler.Me)
		}

		// Protected Routes
		protected := api.Group("/")
		protected.Use(handler.JWTMiddleware(authService))
		{
			// Kalender Route
			protected.GET("/kalender", kalenderHandler.GetCalendar)

			// Toko Routes
			toko := protected.Group("/toko")
			{
				toko.GET("", tokoHandler.GetAll)
				toko.GET("/:toko_id", tokoHandler.GetByID)
				toko.POST("", handler.AdminOnly(), tokoHandler.Create)
			}

			// Karyawan Routes
			karyawan := protected.Group("/karyawan")
			{
				karyawan.GET("", karyawanHandler.GetAll)
				karyawan.GET("/:karyawan_id", karyawanHandler.GetByID)
				karyawan.POST("", handler.AdminOnly(), karyawanHandler.Create)
				karyawan.PATCH("/:karyawan_id", handler.AdminOnly(), karyawanHandler.Update)
			}

			// Jadwal Libur Routes
			jadwal := protected.Group("/jadwal-libur")
			{
				jadwal.GET("", jadwalHandler.GetAll)
				jadwal.GET("/check", jadwalHandler.CheckAvailability)
				jadwal.POST("", jadwalHandler.CreatePlanned)
				jadwal.POST("/unplanned", handler.AdminOnly(), jadwalHandler.CreateUnplanned)
				jadwal.GET("/:id", jadwalHandler.GetByID)
				jadwal.PATCH("/:id", jadwalHandler.Update)
				jadwal.DELETE("/:id", jadwalHandler.Delete)
			}

			// Backup Assignment Routes
			backup := protected.Group("/backup-assignment")
			{
				backup.POST("", handler.AdminOnly(), backupHandler.Create)
				backup.DELETE("/:id", handler.AdminOnly(), backupHandler.Delete)
			}

			// Metrik Routes
			metrik := protected.Group("/metrik")
			metrik.Use(handler.AdminOnly())
			{
				metrik.GET("/karyawan/:karyawan_id", metrikHandler.GetKaryawanMetrik)
				metrik.GET("/toko/:toko_id", metrikHandler.GetTokoMetrik)
			}

			// Konfigurasi Routes
			konfigurasi := protected.Group("/konfigurasi")
			konfigurasi.Use(handler.AdminOnly())
			{
				konfigurasi.GET("", configHandler.GetAll)
				konfigurasi.PATCH("/:key", configHandler.Update)
			}
		}

		// Ping
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	log.Println("Server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
