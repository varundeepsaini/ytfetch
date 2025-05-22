package handler

import (
	"net/http"
	"strconv"

	"ytfetch/internal/service"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	service *service.YouTubeService
}

func NewVideoHandler(service *service.YouTubeService) *VideoHandler {
	return &VideoHandler{
		service: service,
	}
}

// GetLatestVideos handles GET /api/videos
// Query params:
// - cursor: string (optional) - cursor for pagination
// - limit: int (optional, default: 10) - number of videos per page
func (h *VideoHandler) GetLatestVideos(c *gin.Context) {
	cursor := c.Query("cursor")
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Get videos from service
	videos, err := h.service.GetLatestVideos(c.Request.Context(), cursor, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch videos",
		})
		return
	}

	c.JSON(http.StatusOK, videos)
}
