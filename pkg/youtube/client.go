package youtube

import (
	"context"
	"fmt"
	"time"
	"ytfetch/internal/models"

	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	service    *youtube.Service
	apiKeys    []string
	currentKey int
	logger     *zap.Logger
}

const (
	quotaExceeded = "googleapi: Error 403: The request cannot be completed because you have exceeded your <a href=\"/youtube/v3/getting-started#quota\">quota</a>., quotaExceeded"
)

func NewClient(apiKeys []string) (*Client, error) {
	if len(apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided")
	}

	logger, _ := zap.NewProduction()

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKeys[0]))
	if err != nil {
		logger.Error("Failed to create YouTube service",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create YouTube service: %v", err)
	}

	return &Client{
		service:    service,
		apiKeys:    apiKeys,
		currentKey: 0,
		logger:     logger,
	}, nil
}

func (c *Client) FetchLatestVideos(ctx context.Context, query string, publishedAfter string) ([]models.Video, error) {
	c.logger.Info("Fetching latest videos from YouTube API",
		zap.String("query", query),
		zap.String("published_after", publishedAfter),
	)

	call := c.service.Search.List([]string{"snippet"}).
		Q(query).
		MaxResults(25).
		Order("date").
		Type("video").
		PublishedAfter(publishedAfter)

	response, err := call.Do()
	if err != nil {
		if err.Error() == quotaExceeded {
			c.logger.Warn("Quota exceeded, rotating API key")
			if err := c.rotateKey(); err != nil {
				c.logger.Error("Failed to rotate API key",
					zap.Error(err),
				)
				return nil, err
			}
			return c.FetchLatestVideos(ctx, query, publishedAfter)
		}
		c.logger.Error("Failed to fetch videos",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to fetch videos: %v", err)
	}

	videos := make([]models.Video, 0, len(response.Items))
	for _, item := range response.Items {
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			c.logger.Error("Failed to parse published date",
				zap.Error(err),
				zap.String("published_at", item.Snippet.PublishedAt),
			)
			return nil, fmt.Errorf("failed to parse published date: %v", err)
		}

		video := models.Video{
			ID:           item.Id.VideoId,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishedAt:  publishedAt,
			ThumbnailURL: item.Snippet.Thumbnails.Default.Url,
			ChannelTitle: item.Snippet.ChannelTitle,
			ChannelID:    item.Snippet.ChannelId,
		}
		videos = append(videos, video)
	}

	c.logger.Info("Successfully fetched videos",
		zap.Int("count", len(videos)),
	)
	return videos, nil
}

func (c *Client) rotateKey() error {
	c.currentKey = (c.currentKey + 1) % len(c.apiKeys)

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(c.apiKeys[c.currentKey]))
	if err != nil {
		c.logger.Error("Failed to create YouTube service with new key",
			zap.Error(err),
			zap.Int("key_index", c.currentKey),
		)
		return fmt.Errorf("failed to create YouTube service with new key: %v", err)
	}

	c.service = service
	c.logger.Info("Rotated to new API key",
		zap.Int("key_index", c.currentKey),
	)
	return nil
}
