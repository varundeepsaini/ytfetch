package models

import (
	"time"
)

type Video struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title" gorm:"type:varchar(255);not null;index:idx_title"`
	Description  string    `json:"description" gorm:"type:text"`
	PublishedAt  time.Time `json:"published_at" gorm:"index:idx_published_at;not null"`
	ThumbnailURL string    `json:"thumbnail_url" gorm:"type:varchar(255)"`
	ChannelTitle string    `json:"channel_title" gorm:"type:varchar(255);index:idx_channel"`
	ChannelID    string    `json:"channel_id" gorm:"type:varchar(255);index:idx_channel"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime;index:idx_created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Video) TableName() string {
	return "videos"
}

type VideoResponse struct {
	Videos     []Video `json:"videos"`
	Total      int64   `json:"total"`
	Limit      int     `json:"limit"`
	NextCursor string  `json:"next_cursor,omitempty"`
	HasMore    bool    `json:"has_more"`
}
