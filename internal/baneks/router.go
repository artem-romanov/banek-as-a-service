package baneks

import (
	"baneks.com/internal/baneks/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func InitBanekRouter(group *echo.Group) *echo.Group {
	mainGroup := group.Group("/baneks")

	// Adding rate limiter to avoid hitting banek servers too hard
	mainGroup.Use(
		middleware.RateLimiter(
			middleware.NewRateLimiterMemoryStore(rate.Limit(5)),
		),
	)

	// Adding routes
	mainGroup.GET("/random", handlers.GetRandomBanek)
	mainGroup.GET("/:slug", handlers.GetBanekBySlug)

	return mainGroup
}
