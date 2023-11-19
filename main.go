// main.go

package main

import (
	"blogbackend/models"
	"blogbackend/routes"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Create a PostgreSQL database connection using gorm
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	var errConnect error
	db, errConnect = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errConnect != nil {
		panic(errConnect)
	}

	// Check the database connection
	sqlDB, errPing := db.DB()
	if errPing != nil {
		panic(errPing)
	}

	fmt.Println("Successfully connected to the database")

	// AutoMigrate will create the 'posts' table if it doesn't exist
	errAutoMigrate := db.AutoMigrate(&models.Post{})
	if errAutoMigrate != nil {
		panic(errAutoMigrate)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
}

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Set up routes
	routes.SetupPostRoutes(router, db)

	// Run the application on port 8080
	router.Run(":8080")
}
