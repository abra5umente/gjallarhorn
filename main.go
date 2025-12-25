package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "gjallarhorn/docs"
)

// @title Gjallarhorn API
// @version 1.0
// @description Uptime monitoring service API
// @host localhost:8080
// @BasePath /api

// @tag.name Services
// @tag.description Service monitoring operations
// @tag.name Bulk Operations
// @tag.description Bulk service operations with all-or-nothing semantics
// @tag.name Notifications
// @tag.description Notification configuration

//go:embed dist/*
var frontendFiles embed.FS

func main() {
	// Load .env file
	godotenv.Load()

	// Initialize Echo
	e := echo.New()

	// Initialize validator
	validator := validator.New()
	e.Validator = &CustomValidator{validator: validator}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Initialize services
	monitorService := NewMonitorService()
	notificationService := NewNotificationService()

	// Start background monitoring
	go monitorService.StartMonitoring(notificationService)

	// API routes
	api := e.Group("/api")
	api.GET("/services", monitorService.GetServices)
	api.POST("/services", monitorService.CreateService)
	api.PUT("/services/:id", monitorService.UpdateService)
	api.DELETE("/services/:id", monitorService.DeleteService)
	api.GET("/services/:id/status", monitorService.GetServiceStatus)

	// Bulk operations
	api.POST("/services/bulk", monitorService.BulkCreateServices)
	api.PUT("/services/bulk", monitorService.BulkUpdateServices)
	api.DELETE("/services/bulk", monitorService.BulkDeleteServices)

	// Notifications
	api.POST("/notifications/config", notificationService.UpdateConfig)
	api.GET("/notifications/config", func(c echo.Context) error {
		return c.JSON(http.StatusOK, notificationService.GetConfig())
	})

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Serve frontend files with proper MIME types
	e.GET("/*", func(c echo.Context) error {
		path := c.Request().URL.Path

		// Handle root path - serve index.html
		if path == "/" || path == "" {
			path = "/index.html"
		}

		// Remove leading slash for embedded filesystem
		if strings.HasPrefix(path, "/") {
			path = path[1:]
		}

		// Read file from embedded filesystem
		file, err := frontendFiles.Open("dist/" + path)
		if err != nil {
			// If file not found, serve index.html for SPA routing
			file, err = frontendFiles.Open("dist/index.html")
			if err != nil {
				return c.String(http.StatusNotFound, "File not found")
			}
		}
		defer file.Close()

		// Set proper MIME type based on file extension
		contentType := "text/plain"
		if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".js") {
			contentType = "application/javascript"
		} else if strings.HasSuffix(path, ".html") {
			contentType = "text/html"
		} else if strings.HasSuffix(path, ".json") {
			contentType = "application/json"
		} else if strings.HasSuffix(path, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
			contentType = "image/jpeg"
		} else if strings.HasSuffix(path, ".svg") {
			contentType = "image/svg+xml"
		} else if strings.HasSuffix(path, ".ico") {
			contentType = "image/x-icon"
		}

		c.Response().Header().Set("Content-Type", contentType)
		return c.Stream(http.StatusOK, contentType, file)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Gjallarhorn server on port %s", port)
	log.Fatal(e.Start(":" + port))
}

// CustomValidator wraps the go-playground validator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates a struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
