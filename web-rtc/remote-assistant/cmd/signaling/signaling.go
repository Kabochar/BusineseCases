package main

import (
	"os"

	"remote-assistant/internal/config"
	"remote-assistant/internal/router"
	"remote-assistant/internal/server"
)

func main() {
	config.LoadConfig()

	svr := server.NewEngine()
	router.NewRouter(svr)
	svr.Run(os.Getenv("SERVER_ADDR"))
}
