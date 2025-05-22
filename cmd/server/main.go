package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"ytfetch/internal/config"
	"ytfetch/internal/repository"
	"ytfetch/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	cfg := config.NewConfig()

	repo := repository.NewVideoRepository(nil)

	youtubeService, err := service.NewYouTubeService(cfg.YouTubeAPIKeys, repo, cfg.SearchQuery, cfg.FetchInterval)
	if err != nil {
		log.Fatalf("Failed to create YouTube service: %v", err)
	}

	youtubeService.StartBackgroundFetch()
	log.Printf("Started background fetch for query: %s", cfg.SearchQuery)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	youtubeService.StopBackgroundFetch()
	log.Println("Stopped background fetch")
}
