package main

import (
	"hr-system/config"
	"hr-system/internal/handlers"
	"hr-system/internal/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()
	r := gin.Default()

	config.DB.AutoMigrate(
		&models.User{},
		&models.AuditLog{},
		&models.PayCycle{},
		&models.Attendance{},
		&models.Overtime{},
	)

	handlers.RegisterRoutes(r)
	log.Println("🚀 Server running at http://localhost:8080")
	r.Run(":8080")
}
