package main

import (
	server "baneks.com/internal/api"
	"baneks.com/internal/api/baneks"
	"baneks.com/internal/api/memes"
)

func main() {
	server := server.InitializeServer()
	g := server.Group("/api")
	baneks.InitBanekRouter(g)
	memes.InitMemesRouter(g)
	server.Logger.Fatal(server.Start("localhost:8888"))
}
