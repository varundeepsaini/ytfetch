package repository

import (
	"context"
	"time"
	"ytfetch/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type VideoRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	logger, _ := zap.NewProduction()
	return &VideoRepository{
		db:     db,
		logger: logger,
	}
}

func (r *VideoRepository) BatchCreate(ctx context.Context, videos []models.Video) error {
	r.logger.Info("Batch creating videos",
		zap.Int("count", len(videos)),
	)

	if err := r.db.WithContext(ctx).Create(&videos).Error; err != nil {
		r.logger.Error("Failed to batch create videos",
			zap.Error(err),
			zap.Int("count", len(videos)),
		)
		return err
	}
	return nil
}

func (r *VideoRepository) CreateOrUpdate(ctx context.Context, video *models.Video) error {
	r.logger.Info("Creating or updating video",
		zap.String("video_id", video.ID),
		zap.String("title", video.Title),
	)

	if err := r.db.WithContext(ctx).Save(video).Error; err != nil {
		r.logger.Error("Failed to create or update video",
			zap.Error(err),
			zap.String("video_id", video.ID),
		)
		return err
	}
	return nil
}

func (r *VideoRepository) GetLatestVideos(ctx context.Context, cursor string, limit int) (*models.VideoResponse, error) {
	r.logger.Info("Getting latest videos",
		zap.String("cursor", cursor),
		zap.Int("limit", limit),
	)

	var videos []models.Video
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Video{}).Count(&total).Error; err != nil {
		r.logger.Error("Failed to count videos",
			zap.Error(err),
		)
		return nil, err
	}

	query := r.db.WithContext(ctx).Order("published_at DESC").Limit(limit + 1)

	if cursor != "" {
		cursorTime, err := time.Parse(time.RFC3339, cursor)
		if err != nil {
			r.logger.Error("Failed to parse cursor",
				zap.Error(err),
				zap.String("cursor", cursor),
			)
			return nil, err
		}
		query = query.Where("published_at < ?", cursorTime)
	}

	if err := query.Find(&videos).Error; err != nil {
		r.logger.Error("Failed to find videos",
			zap.Error(err),
		)
		return nil, err
	}

	// checking for more
	hasMore := len(videos) > limit
	if hasMore {
		videos = videos[:limit]
	}

	var nextCursor string
	if hasMore && len(videos) > 0 {
		nextCursor = videos[len(videos)-1].PublishedAt.UTC().Format(time.RFC3339)
	}

	r.logger.Info("Retrieved latest videos",
		zap.Int("count", len(videos)),
		zap.Int64("total", total),
		zap.Bool("has_more", hasMore),
	)

	return &models.VideoResponse{
		Videos:     videos,
		Total:      total,
		Limit:      limit,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (r *VideoRepository) GetLatestPublishedAt(ctx context.Context) (string, error) {
	r.logger.Info("Getting latest published at time")

	var video models.Video
	if err := r.db.WithContext(ctx).
		Order("published_at DESC").
		First(&video).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			defaultTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
			r.logger.Info("No videos found, using default time",
				zap.String("default_time", defaultTime),
			)
			return defaultTime, nil
		}
		r.logger.Error("Failed to get latest published at",
			zap.Error(err),
		)
		return "", err
	}

	latestTime := video.PublishedAt.Add(1 * time.Second).Format(time.RFC3339)
	r.logger.Info("Retrieved latest published at time",
		zap.String("latest_time", latestTime),
	)
	return latestTime, nil
}
