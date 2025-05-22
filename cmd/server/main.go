package main

import (
	"os"
	"os/signal"
	"syscall"
	"ytfetch/internal/config"
	"ytfetch/internal/handler"
	"ytfetch/internal/repository"
	"ytfetch/internal/router"
	"ytfetch/internal/service"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Warn("Error loading .env file",
			zap.Error(err),
		)
	}

	cfg := config.NewConfig()

	db, err := config.InitDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database",
			zap.Error(err),
		)
	}

	repo := repository.NewVideoRepository(db)

	youtubeService, err := service.NewYouTubeService(cfg.YouTubeAPIKeys, repo, cfg.SearchQuery, cfg.FetchInterval)
	if err != nil {
		logger.Fatal("Failed to create YouTube service",
			zap.Error(err),
		)
	}

	youtubeService.StartBackgroundFetch()
	logger.Info("Started background fetch",
		zap.String("query", cfg.SearchQuery),
	)

	videoHandler := handler.NewVideoHandler(youtubeService)

	router := router.SetupRouter(videoHandler)

	go func() {
		if err := router.Run(":8080"); err != nil {
			logger.Fatal("Failed to start server",
				zap.Error(err),
			)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	youtubeService.StopBackgroundFetch()
	logger.Info("Stopped background fetch")
}
