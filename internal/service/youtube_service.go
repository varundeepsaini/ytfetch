package service

import (
	"context"
	"fmt"
	"time"

	"ytfetch/internal/models"
	"ytfetch/internal/repository"
	"ytfetch/pkg/youtube"

	"go.uber.org/zap"
)

type YouTubeService struct {
	client        *youtube.Client
	repo          *repository.VideoRepository
	searchQuery   string
	fetchInterval time.Duration
	stopChan      chan struct{}
	logger        *zap.Logger
}

func NewYouTubeService(apiKeys []string, repo *repository.VideoRepository, searchQuery string, fetchInterval time.Duration) (*YouTubeService, error) {
	client, err := youtube.NewClient(apiKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to create YouTube client: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %v", err)
	}

	return &YouTubeService{
		client:        client,
		repo:          repo,
		searchQuery:   searchQuery,
		fetchInterval: fetchInterval,
		stopChan:      make(chan struct{}),
		logger:        logger,
	}, nil
}

func (s *YouTubeService) StartBackgroundFetch() {
	ticker := time.NewTicker(s.fetchInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := s.fetchAndStoreVideos(); err != nil {
					s.logger.Error("Error fetching videos",
						zap.Error(err),
					)
				}
			case <-s.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *YouTubeService) StopBackgroundFetch() {
	s.logger.Info("Stopping background fetch")
	close(s.stopChan)
}

func (s *YouTubeService) GetLatestVideos(ctx context.Context, cursor string, limit int) (*models.VideoResponse, error) {
	return s.repo.GetLatestVideos(ctx, cursor, limit)
}

func (s *YouTubeService) fetchAndStoreVideos() error {
	ctx := context.Background()

	latestPublishedAt, err := s.repo.GetLatestPublishedAt(ctx)
	if err != nil {
		s.logger.Error("Failed to get latest published at",
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("Fetching videos",
		zap.String("query", s.searchQuery),
		zap.String("published_after", latestPublishedAt),
	)

	videos, err := s.client.FetchLatestVideos(ctx, s.searchQuery, latestPublishedAt)
	if err != nil {
		s.logger.Error("Failed to fetch videos",
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("Fetched videos from YouTube API",
		zap.Int("count", len(videos)),
	)

	if len(videos) > 0 {
		if err := s.repo.BatchCreate(ctx, videos); err != nil {
			s.logger.Error("Failed to store videos",
				zap.Error(err),
			)
			return err
		}
		s.logger.Info("Stored new videos",
			zap.Int("count", len(videos)),
		)
	} else {
		s.logger.Info("No new videos found",
			zap.String("query", s.searchQuery),
		)
	}
	return nil
}
