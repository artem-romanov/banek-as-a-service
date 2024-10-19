package memes

import (
	"baneks.com/internal/memes/handlers"
	"github.com/labstack/echo/v4"
)

func InitMemesRouter(group *echo.Group) *echo.Group {
	mainGroup := group.Group("/memes")

	mainGroup.GET("/random", handlers.GetRandomMemes)

	return mainGroup
}
