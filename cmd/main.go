package main

import (
	"log"

	server "baneks.com/internal/api"
	"baneks.com/internal/api/baneks"
	"baneks.com/internal/api/memes"
	"baneks.com/internal/api/middlewares"
	"baneks.com/internal/config"
)

func main() {
	config, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Load config error: %s", err)
		return
	}
	guard := middlewares.New(config.ApiKey)

	server := server.InitializeServer()
	g := server.Group("/api")

	// global middlewares init
	g.Use(guard.GuardWithSecretMiddleware)

	// router init
	baneks.InitBanekRouter(g)
	memes.InitMemesRouter(g)

	// starting the server
	server.Logger.Fatal(server.Start("localhost:8888"))
}
