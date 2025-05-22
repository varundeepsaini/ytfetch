package models

import (
	"time"
)

type Video struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title" gorm:"type:varchar(255);not null"`
	Description  string    `json:"description" gorm:"type:text"`
	PublishedAt  time.Time `json:"published_at" gorm:"index;not null"`
	ThumbnailURL string    `json:"thumbnail_url" gorm:"type:varchar(255)"`
	ChannelTitle string    `json:"channel_title" gorm:"type:varchar(255)"`
	ChannelID    string    `json:"channel_id" gorm:"type:varchar(255)"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type VideoResponse struct {
	Videos []Video `json:"videos"`
	Total  int64   `json:"total"`
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
}
