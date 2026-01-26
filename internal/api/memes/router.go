package memes

import (
	"baneks.com/internal/api/memes/handlers"
	"github.com/labstack/echo/v5"
)

func InitMemesRouter(group *echo.Group) *echo.Group {
	mainGroup := group.Group("/memes")

	mainGroup.GET("/random", handlers.GetRandomMemes)

	return mainGroup
}
