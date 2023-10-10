package router

import (
	"basic-arch/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine) {
	r.GET("/ping", handler.Ping)
}
