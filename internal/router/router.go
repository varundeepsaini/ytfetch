package router

import (
	"ytfetch/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(videoHandler *handler.VideoHandler) *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	videos := v1.Group("/videos")
	videos.GET("", videoHandler.GetLatestVideos)
	return router
}
