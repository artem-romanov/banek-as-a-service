package ai

import (
	"baneks.com/internal/api/ai/handlers"
	aiclient "baneks.com/pkg/ai"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func InitAiRouter(group *echo.Group, aiClient *aiclient.AiClient) *echo.Group {
	mainGroup := group.Group("/ai")

	mainGroup.Use(
		middleware.RateLimiter(
			middleware.NewRateLimiterMemoryStore(3),
		),
	)

	chatHandler := handlers.NewChatWithAiHandler(aiClient)
	mainGroup.POST("/chat", chatHandler.ChatWithAi)

	return mainGroup
}
