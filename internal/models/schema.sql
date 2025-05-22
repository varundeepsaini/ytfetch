CREATE DATABASE IF NOT EXISTS ytfetch;

USE ytfetch;

CREATE TABLE IF NOT EXISTS videos (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    published_at DATETIME NOT NULL,
    thumbnail_url VARCHAR(255),
    channel_title VARCHAR(255),
    channel_id VARCHAR(255),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_published_at (published_at),
    INDEX idx_channel_id (channel_id)
);

CREATE INDEX idx_channel_published ON videos (channel_id, published_at);
CREATE INDEX idx_published_range ON videos (published_at, created_at);
CREATE INDEX idx_title_published ON videos (title, published_at);
CREATE INDEX idx_channel_title ON videos (channel_id, title);