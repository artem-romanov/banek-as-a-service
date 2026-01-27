package baneks

import (
	"baneks.com/internal/api/baneks/handlers"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func InitBanekRouter(group *echo.Group) *echo.Group {
	mainGroup := group.Group("/baneks")

	// Adding rate limiter to avoid hitting banek servers too hard
	mainGroup.Use(
		middleware.RateLimiter(
			middleware.NewRateLimiterMemoryStore(5),
		),
	)

	// Adding routes
	mainGroup.GET("/random", handlers.GetRandomBanek)
	mainGroup.GET("/:slug", handlers.GetBanekBySlug)

	return mainGroup
}
