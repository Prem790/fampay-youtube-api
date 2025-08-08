package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"fampay-youtube-api/internal/services"
	"fampay-youtube-api/internal/utils"
)

type YouTubeSearchHandler struct {
	youtubeService *services.YouTubeService
}

func NewYouTubeSearchHandler(youtubeService *services.YouTubeService) *YouTubeSearchHandler {
	return &YouTubeSearchHandler{
		youtubeService: youtubeService,
	}
}

// LiveSearch - Search YouTube directly for any query (real-time)
func (ysh *YouTubeSearchHandler) LiveSearch(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query cannot be empty",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))
	sortBy := c.DefaultQuery("sort", "latest")

	// Validate parameters
	if page < 1 {
		page = 1
	}
	if pageSize > 50 {
		pageSize = 50
	}
	if pageSize < 1 {
		pageSize = 12
	}

	// Search YouTube directly for any query
	videos, err := ysh.youtubeService.SearchYouTubeLive(query, pageSize, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search YouTube",
			"details": err.Error(),
		})
		return
	}

	// Create response (note: total count is estimated for live search)
	total := int64(len(videos))
	if len(videos) == pageSize {
		total = int64(pageSize * 10) // Estimate more results available
	}

	response := utils.NewPaginatedResponse(videos, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}
