package main

import (
	"log"

	"backend-liburwan/internal/config"
	"backend-liburwan/internal/handler"
	"backend-liburwan/internal/repository"
	"backend-liburwan/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
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

	// Services
	tokoService := service.NewTokoService(tokoRepo)
	karyawanService := service.NewKaryawanService(karyawanRepo)
	jadwalService := service.NewJadwalLiburService(jadwalRepo, karyawanRepo)

	// Handlers
	tokoHandler := handler.NewTokoHandler(tokoService)
	karyawanHandler := handler.NewKaryawanHandler(karyawanService)
	jadwalHandler := handler.NewJadwalLiburHandler(jadwalService)

	// Init Gin Router
	r := gin.Default()

	// API Routes
	api := r.Group("/api/v1")
	{
		// Toko Routes
		toko := api.Group("/toko")
		{
			toko.GET("", tokoHandler.GetAll)
			toko.GET("/:toko_id", tokoHandler.GetByID)
			toko.POST("", handler.AdminOnly(), tokoHandler.Create)
		}

		// Karyawan Routes
		karyawan := api.Group("/karyawan")
		{
			karyawan.GET("", karyawanHandler.GetAll)
			karyawan.GET("/:karyawan_id", karyawanHandler.GetByID)
			karyawan.POST("", handler.AdminOnly(), karyawanHandler.Create)
			karyawan.PATCH("/:karyawan_id", handler.AdminOnly(), karyawanHandler.Update)
		}

		// Jadwal Libur Routes
		jadwal := api.Group("/jadwal-libur")
		{
			jadwal.GET("", jadwalHandler.GetAll)
			jadwal.GET("/check", jadwalHandler.CheckAvailability)
			jadwal.POST("", jadwalHandler.CreatePlanned)
			jadwal.POST("/unplanned", handler.AdminOnly(), jadwalHandler.CreateUnplanned)
			jadwal.GET("/:id", jadwalHandler.GetByID)
			jadwal.PATCH("/:id", jadwalHandler.Update)
			jadwal.DELETE("/:id", jadwalHandler.Delete)
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
