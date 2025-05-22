package router

import (
	"ytfetch/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(videoHandler *handler.VideoHandler) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := router.Group("/api/v1")
	videos := api.Group("/videos")
	videos.GET("", videoHandler.GetLatestVideos)

	return router
}
