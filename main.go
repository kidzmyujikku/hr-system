package main

import (
	"hr-system/config"
	"hr-system/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	r := gin.Default()

	handlers.RegisterRoutes(r)
	log.Println("🚀 Server running at http://localhost:8080")
	r.Run(":8080")
}