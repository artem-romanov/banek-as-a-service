package main

import (
	"baneks.com/internal/baneks"
	"baneks.com/internal/memes"
	"baneks.com/internal/server"
)

func main() {
	server := server.InitializeServer()
	g := server.Group("/api")
	baneks.InitBanekRouter(g)
	memes.InitMemesRouter(g)
	server.Logger.Fatal(server.Start("localhost:8888"))
}
