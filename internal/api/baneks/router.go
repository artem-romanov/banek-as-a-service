package baneks

import (
	"net/http"

	"baneks.com/internal/api/baneks/handlers"
	"baneks.com/internal/loaders/banekloader"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func InitBanekRouter(group *echo.Group, httpClient *http.Client) *echo.Group {
	mainGroup := group.Group("/baneks")

	// Adding rate limiter to avoid hitting banek servers too hard
	mainGroup.Use(
		middleware.RateLimiter(
			middleware.NewRateLimiterMemoryStore(5),
		),
	)

	banekBalancer := banekloader.NewBanekBalancer(httpClient)
	banekSiteLoader := banekloader.NewBaneksSiteLoader(httpClient)
	h := handlers.NewHandlers(banekBalancer, banekSiteLoader)

	// Adding routes
	mainGroup.GET("/random", h.GetRandomBanek)
	mainGroup.GET("/:slug", h.GetBanekBySlug)

	return mainGroup
}
