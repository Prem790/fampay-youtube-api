package worker

import (
	"log"
	"time"

	"fampay-youtube-api/internal/config"
	"fampay-youtube-api/internal/repository"
	"fampay-youtube-api/internal/services"
)

type VideoFetcher struct {
	videoRepo      *repository.VideoRepository
	youtubeService *services.YouTubeService
	fetchInterval  time.Duration
	stopChan       chan struct{}
	config         config.YouTubeConfig
}

func NewVideoFetcher(videoRepo *repository.VideoRepository, youtubeConfig config.YouTubeConfig) *VideoFetcher {
	youtubeService := services.NewYouTubeService(
		youtubeConfig.APIKeys,
		youtubeConfig.SearchQueries,
		youtubeConfig.MaxResultsPerQuery,
		youtubeConfig.RegionCode,
		youtubeConfig.RelevanceLanguage,
	)

	return &VideoFetcher{
		videoRepo:      videoRepo,
		youtubeService: youtubeService,
		fetchInterval:  time.Duration(youtubeConfig.FetchInterval) * time.Second, // EXACTLY as per config (10 seconds)
		stopChan:       make(chan struct{}),
		config:         youtubeConfig,
	}
}

func (vf *VideoFetcher) Start() {
	log.Printf("🚀 Starting video fetcher (FamPay Requirements Compliance):")
	log.Printf("📋 Search queries: %v", vf.config.SearchQueries)
	log.Printf("🔑 API keys: %d keys available", len(vf.config.APIKeys))
	log.Printf("⏰ Fetch interval: %d seconds (as per requirements)", vf.config.FetchInterval)
	log.Printf("📊 Max results per query: %d", vf.config.MaxResultsPerQuery)
	log.Printf("🌍 Region: %s, Language: %s", vf.config.RegionCode, vf.config.RelevanceLanguage)

	// Initial fetch
	vf.fetchAndStore()

	// FamPay Requirement: Continuous background fetching at specified interval
	ticker := time.NewTicker(vf.fetchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			vf.fetchAndStore()
		case <-vf.stopChan:
			log.Println("Video fetcher stopped")
			return
		}
	}
}

func (vf *VideoFetcher) Stop() {
	close(vf.stopChan)
}

func (vf *VideoFetcher) fetchAndStore() {
	startTime := time.Now()

	// Get last published date from database to avoid duplicates
	lastVideo, err := vf.videoRepo.GetLatest()

	// FamPay Requirement: Fetch latest videos - use reasonable time window
	publishedAfter := time.Now().Add(-2 * time.Hour)
	if err == nil && lastVideo != nil && !lastVideo.PublishedAt.IsZero() {
		// Use last published date but ensure we don't miss videos
		publishedAfter = lastVideo.PublishedAt.Add(-5 * time.Minute)
	}

	log.Printf("🔍 Fetching latest videos published after: %s", publishedAfter.Format("2006-01-02 15:04:05"))

	// FamPay Requirement: Fetch for ALL predefined search queries
	videos, err := vf.youtubeService.FetchLatestVideosForAllQueries(publishedAfter)
	if err != nil {
		log.Printf("❌ Error fetching videos: %v", err)

		// Log API key status for debugging - FamPay Bonus: Multiple API key support
		status := vf.youtubeService.GetAPIKeyStatus()
		log.Printf("🔑 API Status: %d/%d keys working, next retry in %v",
			status["working_keys"], status["total_keys"], vf.fetchInterval)
		return
	}

	if len(videos) == 0 {
		log.Printf("📭 No new videos found (search completed in %v)", time.Since(startTime))
		return
	}

	// FamPay Requirement: Store video data in database
	stored := 0
	skipped := 0
	errors := 0

	for _, video := range videos {
		// Check if video already exists to avoid duplicates
		existingVideo, err := vf.videoRepo.GetByVideoID(video.VideoID)

		if err != nil {
			log.Printf("⚠️ Error checking video %s: %v", video.VideoID, err)
			errors++
			continue
		}

		if existingVideo == nil {
			// FamPay Requirement: Store video with all required fields
			if err := vf.videoRepo.Create(video); err != nil {
				log.Printf("❌ Error storing video %s: %v", video.VideoID, err)
				errors++
				continue
			}
			stored++
		} else {
			skipped++
		}
	}

	duration := time.Since(startTime)
	log.Printf("✅ Fetch cycle completed: %d stored, %d duplicates skipped, %d errors (took %v)",
		stored, skipped, errors, duration)

	// Log API status for FamPay Bonus: Multiple API key management
	if stored > 0 || errors > 0 {
		status := vf.youtubeService.GetAPIKeyStatus()
		log.Printf("🔑 API Key Status: Using key %v, %d/%d keys working",
			status["current_key_index"], status["working_keys"], status["total_keys"])
	}
}
