package repository

import (
	"context"
	"fmt"
	"time"
	"ytfetch/internal/models"

	"gorm.io/gorm"
)

// VideoRepository handles database operations for videos
type VideoRepository struct {
	db *gorm.DB
}

// NewVideoRepository creates a new VideoRepository instance
func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

// CreateOrUpdate prints the video data instead of saving to database
func (r *VideoRepository) CreateOrUpdate(ctx context.Context, video *models.Video) error {
	fmt.Printf("\n=== New Video ===\n")
	fmt.Printf("ID: %s\n", video.ID)
	fmt.Printf("Title: %s\n", video.Title)
	fmt.Printf("Description: %s\n", video.Description)
	fmt.Printf("Published At: %s\n", video.PublishedAt.Format(time.RFC3339))
	fmt.Printf("Thumbnail URL: %s\n", video.ThumbnailURL)
	fmt.Printf("Channel Title: %s\n", video.ChannelTitle)
	fmt.Printf("Channel ID: %s\n", video.ChannelID)
	fmt.Printf("================\n\n")
	return nil
}

// GetLatestVideos returns paginated videos sorted by published date
func (r *VideoRepository) GetLatestVideos(ctx context.Context, page, limit int) (*models.VideoResponse, error) {
	// For testing, return empty response
	return &models.VideoResponse{
		Videos: []models.Video{},
		Total:  0,
		Page:   page,
		Limit:  limit,
	}, nil
}

// GetLatestPublishedAt returns the latest published_at timestamp from the database
func (r *VideoRepository) GetLatestPublishedAt(ctx context.Context) (time.Time, error) {
	// For testing, return timestamp 24 hours ago
	return time.Now().Add(-24 * time.Hour), nil
}
