package router

import (
	"github.com/gin-gonic/gin"
	"github.com/org/task-api/internal/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures the Gin router with all routes
func SetupRouter(taskHandler *handlers.TaskHandler) *gin.Engine {
	r := gin.Default()

	// Add CORS middleware
	r.Use(corsMiddleware())

	// API routes
	api := r.Group("/api")
	{
		task := api.Group("/task")
		{
			task.POST("", taskHandler.AddTask)
			task.GET("", taskHandler.GetTasks)
			task.PUT("", taskHandler.UpdateTask)
			task.GET("/:id", taskHandler.GetTask)
			task.PUT("/:id/complete", taskHandler.CompleteTask)
			task.DELETE("/:id", taskHandler.DeleteTask)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}