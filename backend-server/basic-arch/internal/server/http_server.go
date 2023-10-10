package server

import (
	"os"

	"github.com/gin-gonic/gin"
)

func NewEngine() *gin.Engine {
	// Set gin mode
	gin.SetMode(os.Getenv("SERVER_MODE"))

	r := gin.Default()
	return r
}
