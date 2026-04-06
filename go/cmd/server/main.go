package main

import (
	"fmt"
	"log"
	"os"

	"github.com/org/task-api/internal/handlers"
	"github.com/org/task-api/internal/models"
	"github.com/org/task-api/internal/repository"
	"github.com/org/task-api/internal/router"
	"github.com/org/task-api/internal/service"
	"github.com/org/task-api/internal/validator"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// @title Task Management API
// @version 1.0
// @description REST API for managing tasks
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

func main() {
	// Database configuration
	dsn := getDSN()
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize components
	taskValidator := validator.NewTaskValidator()
	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo, taskValidator)
	taskHandler := handlers.NewTaskHandler(taskService)

	// Setup router
	r := router.SetupRouter(taskHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getDSN() string {
	user := getEnv("DB_USER", "secretuser")
	password := getEnv("DB_PASSWORD", "thisisasupersecretpassworddontyouthink")
	host := getEnv("DB_SERVER", "mysql")
	port := getEnv("DB_PORT", "3306")
	dbname := getEnv("DB_NAME", "task-list")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}