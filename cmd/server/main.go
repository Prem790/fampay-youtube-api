package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"fampay-youtube-api/internal/api/routes"
	"fampay-youtube-api/internal/config"
	"fampay-youtube-api/internal/repository"
	"fampay-youtube-api/internal/worker"
	"fampay-youtube-api/pkg/database"
	"fampay-youtube-api/pkg/redis"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize MongoDB
	db, err := database.NewMongoDB(cfg.MongoDB)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Initialize Redis
	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// Initialize repository
	videoRepo := repository.NewVideoRepository(db)

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize router
	router := routes.SetupRouter(videoRepo, redisClient, cfg)

	// Start background worker
	videoFetcher := worker.NewVideoFetcher(videoRepo, cfg.YouTube)
	go videoFetcher.Start()

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	fmt.Printf("Server started on :%s\n", cfg.Server.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
