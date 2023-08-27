package router

import (
	"github.com/gin-gonic/gin"

	"remote-assistant/internal/handler"
)

func NewRouter(r *gin.Engine) {
	r.GET("/ping", handler.Ping)

	signaling := r.Group("/signaling")
	{
		signaling.GET("", handler.SignalingServer)
		signaling.GET("/server/info", handler.ServerInfo)
	}
}
