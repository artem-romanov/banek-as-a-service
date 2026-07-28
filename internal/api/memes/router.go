package memes

import (
	"baneks.com/internal/api/memes/handlers"
	memesloader "baneks.com/internal/loaders/memes_loader"
	"github.com/labstack/echo/v5"
)

func InitMemesRouter(group *echo.Group, memeLoader memesloader.MemeLoader) *echo.Group {
	mainGroup := group.Group("/memes")

	h := handlers.NewHandlers(memeLoader)
	mainGroup.GET("/random", h.GetRandomMemes)

	return mainGroup
}
