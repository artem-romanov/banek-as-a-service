package memegenerator

import (
	"baneks.com/internal/api/meme_generator/handlers"
	"baneks.com/pkg/memer"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func InitMemeGeneratorRouter(group *echo.Group, memer *memer.Memer) *echo.Group {
	mainGroup := group.Group("/meme-generator")

	mainGroup.Use(
		middleware.RateLimiter(
			middleware.NewRateLimiterMemoryStore(1),
		),
	)

	createMemeHandler := handlers.NewCreateMemeHandler(memer)
	mainGroup.POST("", createMemeHandler.CreateMeme)

	return mainGroup
}
