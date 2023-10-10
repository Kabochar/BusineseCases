package main

import (
	"os"

	"basic-arch/internal/cache"
	"basic-arch/internal/config"
	"basic-arch/internal/router"
	"basic-arch/internal/server"
)

func main() {
	config.LoadConfig()
	cache.NewCache()

	svr := server.NewEngine()
	router.NewRouter(svr)
	svr.Run(os.Getenv("SERVER_ADDR"))
}
