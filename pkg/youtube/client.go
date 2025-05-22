package youtube

import (
	"context"
	"fmt"
	"time"

	"ytfetch/internal/models"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Client represents a YouTube API client
type Client struct {
	service    *youtube.Service
	apiKeys    []string
	currentKey int
}

// NewClient creates a new YouTube API client
func NewClient(apiKeys []string) (*Client, error) {
	if len(apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided")
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKeys[0]))
	if err != nil {
		return nil, fmt.Errorf("failed to create YouTube service: %v", err)
	}

	return &Client{
		service:    service,
		apiKeys:    apiKeys,
		currentKey: 0,
	}, nil
}

func (c *Client) FetchLatestVideos(ctx context.Context, query string, publishedAfter time.Time) ([]models.Video, error) {
	call := c.service.Search.List([]string{"snippet"}).
		Q(query).
		MaxResults(50).
		Order("date").
		Type("video").
		PublishedAfter(publishedAfter.Format(time.RFC3339))

	response, err := call.Do()
	if err != nil {
		// If quota exceeded, try next API key
		if err.Error() == "quotaExceeded" {
			if err := c.rotateKey(); err != nil {
				return nil, err
			}
			return c.FetchLatestVideos(ctx, query, publishedAfter)
		}
		return nil, fmt.Errorf("failed to fetch videos: %v", err)
	}

	videos := make([]models.Video, 0, len(response.Items))
	for _, item := range response.Items {
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
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

	return videos, nil
}

func (c *Client) rotateKey() error {
	c.currentKey = (c.currentKey + 1) % len(c.apiKeys)

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(c.apiKeys[c.currentKey]))
	if err != nil {
		return fmt.Errorf("failed to create YouTube service with new key: %v", err)
	}

	c.service = service
	return nil
}
