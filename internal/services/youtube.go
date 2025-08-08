package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"fampay-youtube-api/internal/models"
)

type YouTubeService struct {
	apiKeys            []string
	currentKeyIdx      int
	mutex              sync.RWMutex
	searchQueries      []string
	maxResultsPerQuery int
	keyQuotaStatus     map[int]time.Time
	regionCode         string
	relevanceLanguage  string
}

func NewYouTubeService(apiKeys []string, searchQueries []string, maxResultsPerQuery int, regionCode, relevanceLanguage string) *YouTubeService {
	log.Printf("Initializing YouTube service with %d API keys and %d search queries", len(apiKeys), len(searchQueries))
	return &YouTubeService{
		apiKeys:            apiKeys,
		searchQueries:      searchQueries,
		maxResultsPerQuery: maxResultsPerQuery,
		keyQuotaStatus:     make(map[int]time.Time),
		regionCode:         regionCode,
		relevanceLanguage:  relevanceLanguage,
	}
}

// FamPay Requirement: Fetch latest videos for predefined search queries
func (ys *YouTubeService) FetchLatestVideosForAllQueries(publishedAfter time.Time) ([]*models.Video, error) {
	var allVideos []*models.Video

	// FamPay Requirement: Fetch for ALL predefined search queries
	for _, query := range ys.searchQueries {
		log.Printf("🔍 Fetching latest videos for query: '%s'", query)

		videos, err := ys.fetchLatestVideosForQuery(query, publishedAfter)
		if err != nil {
			log.Printf("❌ Error fetching videos for query '%s': %v", query, err)
			// Continue with other queries even if one fails
			continue
		}

		allVideos = append(allVideos, videos...)
		log.Printf("✅ Found %d videos for '%s'", len(videos), query)
	}

	log.Printf("📊 Total videos fetched: %d from %d queries", len(allVideos), len(ys.searchQueries))
	return allVideos, nil
}

func (ys *YouTubeService) fetchLatestVideosForQuery(query string, publishedAfter time.Time) ([]*models.Video, error) {
	service, err := ys.getYouTubeService()
	if err != nil {
		return nil, fmt.Errorf("failed to create YouTube service: %w", err)
	}

	// FamPay Requirement: YouTube API call with proper parameters
	call := service.Search.List([]string{"snippet"}).
		Q(query).
		Type("video").                                       // Only videos
		Order("date").                                       // Latest first (FamPay requirement)
		PublishedAfter(publishedAfter.Format(time.RFC3339)). // Latest videos only
		MaxResults(int64(ys.maxResultsPerQuery)).            // Configurable results
		RegionCode(ys.regionCode).                           // Regional content
		RelevanceLanguage(ys.relevanceLanguage).             // Language preference
		SafeSearch("moderate")                               // Safe content

	startTime := time.Now()
	response, err := call.Do()
	apiCallDuration := time.Since(startTime)

	if err != nil {
		// FamPay Bonus: Multiple API key support - rotate on failure
		log.Printf("🔄 API key %d quota exhausted, rotating to next key...", ys.currentKeyIdx+1)
		ys.markKeyAsFailed(ys.currentKeyIdx)

		if ys.rotateToNextWorkingKey() {
			log.Printf("✅ Retrying with API key %d", ys.currentKeyIdx+1)
			return ys.fetchLatestVideosForQuery(query, publishedAfter)
		}
		return nil, fmt.Errorf("all API keys exhausted: %w", err)
	}

	// FamPay Requirement: Extract and store required video fields
	var videos []*models.Video
	for _, item := range response.Items {
		publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)

		video := &models.Video{
			VideoID:      item.Id.VideoId,           // Required field
			Title:        item.Snippet.Title,        // Required field
			Description:  item.Snippet.Description,  // Required field
			PublishedAt:  publishedAt,               // Required field (for sorting)
			ChannelTitle: item.Snippet.ChannelTitle, // Additional useful field
			ChannelID:    item.Snippet.ChannelId,    // Additional useful field
			SearchQuery:  query,                     // Track which query found this
			ThumbnailURL: models.Thumbnail{ // Required field: thumbnail URLs
				Default: getThumbnailURL(item.Snippet.Thumbnails.Default),
				Medium:  getThumbnailURL(item.Snippet.Thumbnails.Medium),
				High:    getThumbnailURL(item.Snippet.Thumbnails.High),
			},
		}
		videos = append(videos, video)
	}

	log.Printf("📹 Fetched %d videos for '%s' in %v (API quota: 100 units used, Key: %d)",
		len(videos), query, apiCallDuration, ys.currentKeyIdx+1)
	return videos, nil
}

// FamPay Bonus: Multiple API key support
func (ys *YouTubeService) markKeyAsFailed(keyIndex int) {
	ys.mutex.Lock()
	defer ys.mutex.Unlock()
	ys.keyQuotaStatus[keyIndex] = time.Now()
	log.Printf("🔑 API key %d marked as exhausted (will reset in ~24 hours)", keyIndex+1)
}

func (ys *YouTubeService) rotateToNextWorkingKey() bool {
	ys.mutex.Lock()
	defer ys.mutex.Unlock()

	originalIdx := ys.currentKeyIdx
	attemptsRemaining := len(ys.apiKeys)

	for attemptsRemaining > 0 {
		ys.currentKeyIdx = (ys.currentKeyIdx + 1) % len(ys.apiKeys)
		attemptsRemaining--

		// Check if this key failed recently (within last 23 hours)
		if failTime, exists := ys.keyQuotaStatus[ys.currentKeyIdx]; exists {
			timeUntilReset := 23*time.Hour - time.Since(failTime)
			if timeUntilReset > 0 {
				log.Printf("⏭️ Skipping API key %d (quota resets in %v)", ys.currentKeyIdx+1, timeUntilReset.Round(time.Hour))
				continue
			} else {
				// Key should be reset now, remove from failed keys
				delete(ys.keyQuotaStatus, ys.currentKeyIdx)
				log.Printf("🔄 API key %d quota should be reset, trying again", ys.currentKeyIdx+1)
			}
		}

		if ys.currentKeyIdx != originalIdx {
			log.Printf("✅ Rotated to API key %d", ys.currentKeyIdx+1)
			return true
		}

		break
	}

	log.Printf("❌ All %d API keys are quota-exhausted. Will retry when quotas reset.", len(ys.apiKeys))
	return false
}

// For search functionality (bonus feature)
func (ys *YouTubeService) SearchYouTubeLive(query string, maxResults int, sortBy string) ([]*models.Video, error) {
	service, err := ys.getYouTubeService()
	if err != nil {
		return nil, fmt.Errorf("failed to create YouTube service: %w", err)
	}

	// Convert sortBy to YouTube API order - Fixed the empty order issue
	var order string
	switch sortBy {
	case "oldest":
		order = "date"
	case "title":
		order = "title"
	case "channel":
		order = "relevance"
	case "latest":
		order = "date"
	default:
		order = "relevance" // Safe default
	}

	call := service.Search.List([]string{"snippet"}).
		Q(query).
		Type("video").
		Order(order).
		MaxResults(int64(maxResults)).
		RegionCode(ys.regionCode).
		RelevanceLanguage(ys.relevanceLanguage).
		SafeSearch("moderate")

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("YouTube search failed: %w", err)
	}

	var videos []*models.Video
	for _, item := range response.Items {
		publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)

		video := &models.Video{
			VideoID:      item.Id.VideoId,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishedAt:  publishedAt,
			ChannelTitle: item.Snippet.ChannelTitle,
			ChannelID:    item.Snippet.ChannelId,
			SearchQuery:  query,
			ThumbnailURL: models.Thumbnail{
				Default: getThumbnailURL(item.Snippet.Thumbnails.Default),
				Medium:  getThumbnailURL(item.Snippet.Thumbnails.Medium),
				High:    getThumbnailURL(item.Snippet.Thumbnails.High),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (ys *YouTubeService) GetAPIKeyStatus() map[string]interface{} {
	ys.mutex.RLock()
	defer ys.mutex.RUnlock()

	workingKeys := len(ys.apiKeys) - len(ys.keyQuotaStatus)

	status := make(map[string]interface{})
	status["current_key_index"] = ys.currentKeyIdx + 1
	status["total_keys"] = len(ys.apiKeys)
	status["working_keys"] = workingKeys
	status["failed_keys"] = len(ys.keyQuotaStatus)

	return status
}

func (ys *YouTubeService) getYouTubeService() (*youtube.Service, error) {
	ys.mutex.RLock()
	apiKey := ys.apiKeys[ys.currentKeyIdx]
	ys.mutex.RUnlock()

	ctx := context.Background()
	return youtube.NewService(ctx, option.WithAPIKey(apiKey))
}

func (ys *YouTubeService) GetSearchQueries() []string {
	return ys.searchQueries
}

func getThumbnailURL(thumb *youtube.Thumbnail) string {
	if thumb != nil {
		return thumb.Url
	}
	return ""
}
