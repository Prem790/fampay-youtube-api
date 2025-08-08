package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"fampay-youtube-api/internal/repository"
	"fampay-youtube-api/internal/utils"
)

type VideoHandler struct {
	videoRepo *repository.VideoRepository
}

func NewVideoHandler(videoRepo *repository.VideoRepository) *VideoHandler {
	return &VideoHandler{
		videoRepo: videoRepo,
	}
}

// Enhanced GetVideos with sorting support
func (vh *VideoHandler) GetVideos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))
	sortBy := c.DefaultQuery("sort", "latest") // latest, oldest, title, channel

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

	// Get videos from repository with sorting
	videos, total, err := vh.videoRepo.GetPaginated(page, pageSize, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch videos",
			"details": err.Error(),
		})
		return
	}

	// Create paginated response
	response := utils.NewPaginatedResponse(videos, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}
