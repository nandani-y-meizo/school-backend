package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	// "fitpro/middleware"
	"shared/infra/db/mdb"
	"shared/middleware"
	"shared/pkgs/jwtmanager"

	"github.com/nandani-y-meizo/school-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Connect to MongoDB
	if err := mdb.InitMongo(); err != nil {
		log.Fatalf("Mongo connection failed: %v", err)
	}
	fmt.Println("MongoDB connected")

	// Initialize Vault JWT
	if err := jwtmanager.InitVaultJWT(); err != nil {
		log.Fatalf("Failed to initialize Vault JWT: %v", err)
	}
	fmt.Println("Vault JWT initialized")

	// Create main app router
	app := gin.Default()

	// Enable CORS globally
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-KEY"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	// Logging middleware
	app.Use(middleware.RequestLogger())

	// API group
	api := app.Group("/api/v1")
	api.Use(middleware.AuthcMiddleware())

	// Public API group (no authentication required)
	publicApi := app.Group("/api/v1")

	// Load routes
	routes.Routes(api)
	routes.PublicRoutes(publicApi)

	// Start server
	server := &http.Server{
		Addr:    ":8085",
		Handler: app,
	}

	fmt.Println("🚀 Server running on :8085")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
