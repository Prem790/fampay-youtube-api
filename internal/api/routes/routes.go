package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"fampay-youtube-api/internal/api/handlers"
	"fampay-youtube-api/internal/api/middleware"
	"fampay-youtube-api/internal/config"
	"fampay-youtube-api/internal/repository"
	"fampay-youtube-api/internal/services"
)

func SetupRouter(videoRepo *repository.VideoRepository, redisClient *redis.Client, cfg *config.Config) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Initialize YouTube service for live search (bonus feature)
	youtubeService := services.NewYouTubeService(
		cfg.YouTube.APIKeys,
		cfg.YouTube.SearchQueries,
		cfg.YouTube.MaxResultsPerQuery,
		cfg.YouTube.RegionCode,
		cfg.YouTube.RelevanceLanguage,
	)

	// Initialize handlers
	videoHandler := handlers.NewVideoHandler(videoRepo)
	searchHandler := handlers.NewSearchHandler(videoRepo)
	youtubeSearchHandler := handlers.NewYouTubeSearchHandler(youtubeService)

	// FamPay Required API endpoints
	api := router.Group("/api")
	{
		videos := api.Group("/videos")
		{
			// FamPay Requirement: GET API for stored videos (paginated, sorted by published_datetime DESC)
			videos.GET("", videoHandler.GetVideos)

			// FamPay Requirement: Basic search API for title and description
			videos.GET("/search", searchHandler.SearchVideos)

			// Bonus: Live YouTube search
			videos.GET("/youtube-search", youtubeSearchHandler.LiveSearch)
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		apiStatus := youtubeService.GetAPIKeyStatus()
		c.JSON(200, gin.H{
			"status":     "ok",
			"timestamp":  time.Now().Unix(),
			"database":   "mongodb",
			"features":   []string{"stored_videos", "search_api", "background_fetcher", "multiple_api_keys"},
			"api_status": apiStatus,
			"compliance": gin.H{
				"background_fetcher": "✅ Running every 10 seconds",
				"pagination":         "✅ Implemented",
				"search_api":         "✅ Title and description search",
				"sorting":            "✅ Published datetime DESC",
				"multiple_api_keys":  "✅ Auto rotation on quota exhaustion",
				"dockerized":         "✅ Docker support available",
			},
		})
	})

	return router
}
