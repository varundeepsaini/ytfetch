package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"ytfetch/internal/models"
	"ytfetch/internal/repository"
	"ytfetch/pkg/youtube"
)

// YouTubeService handles the business logic for YouTube video operations
type YouTubeService struct {
	client   *youtube.Client
	repo     *repository.VideoRepository
	query    string
	interval time.Duration
	stopChan chan struct{}
}

// NewYouTubeService creates a new YouTube service
func NewYouTubeService(apiKeys []string, repo *repository.VideoRepository, query string, interval time.Duration) (*YouTubeService, error) {
	client, err := youtube.NewClient(apiKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to create YouTube client: %v", err)
	}

	return &YouTubeService{
		client:   client,
		repo:     repo,
		query:    query,
		interval: interval,
		stopChan: make(chan struct{}),
	}, nil
}

// StartBackgroundFetch starts the background process to fetch videos
func (s *YouTubeService) StartBackgroundFetch() {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := s.fetchAndStoreVideos(); err != nil {
					log.Printf("Error fetching videos: %v", err)
				}
			case <-s.stopChan:
				return
			}
		}
	}()
}

func (s *YouTubeService) StopBackgroundFetch() {
	close(s.stopChan)
}

func (s *YouTubeService) GetLatestVideos(ctx context.Context, page, limit int) (*models.VideoResponse, error) {
	return s.repo.GetLatestVideos(ctx, page, limit)
}

func (s *YouTubeService) fetchAndStoreVideos() error {
	ctx := context.Background()

	publishedAfter, err := s.repo.GetLatestPublishedAt(ctx)
	if err != nil {
		return err
	}

	videos, err := s.client.FetchLatestVideos(ctx, s.query, publishedAfter)
	if err != nil {
		return err
	}

	for _, video := range videos {
		if err := s.repo.CreateOrUpdate(ctx, &video); err != nil {
			log.Printf("Error storing video %s: %v", video.ID, err)
			continue
		}
	}

	return nil
}
