package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"type:varchar(10);not null;check:role IN ('admin','employee')"`
	Salary    int64  `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func hashPassword(pw string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hashed)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file:", err)
	}

	// Get env values
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Format DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbname, port,
	)

	fmt.Println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// Auto migrate the user model
	db.AutoMigrate(&User{})

	// Create 1 admin
	admin := User{
		Username:  "admin",
		Password:  hashPassword("admin123"),
		Role:      "admin",
		Salary:    0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.FirstOrCreate(&admin, User{Username: "admin"})

	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 1; i++ {
		user := User{
			Username:  fmt.Sprintf("employee%d", i),
			Password:  hashPassword("password123"),
			Role:      "employee",
			Salary:    int64(3000000 + rand.Intn(4000000)), // Rp3M–7M
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		db.FirstOrCreate(&user, User{Username: user.Username})
	}

	fmt.Println("✅ Seeder completed successfully")
}
