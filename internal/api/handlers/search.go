package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"fampay-youtube-api/internal/repository"
	"fampay-youtube-api/internal/utils"
)

type SearchHandler struct {
	videoRepo *repository.VideoRepository
}

func NewSearchHandler(videoRepo *repository.VideoRepository) *SearchHandler {
	return &SearchHandler{
		videoRepo: videoRepo,
	}
}

// Enhanced SearchVideos with better partial matching and sorting
func (sh *SearchHandler) SearchVideos(c *gin.Context) {
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

	// Validate sort parameter
	validSorts := map[string]bool{
		"latest": true, "oldest": true, "title": true, "channel": true,
	}
	if !validSorts[sortBy] {
		sortBy = "latest"
	}

	// Search videos using enhanced partial matching
	videos, total, err := sh.videoRepo.Search(query, page, pageSize, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search videos",
			"details": err.Error(),
		})
		return
	}

	// Create paginated response
	response := utils.NewPaginatedResponse(videos, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}
