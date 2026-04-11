package handlers

import (
	"errors"
	"net/http"

	customerrors "baneks.com/internal/custom_errors"
	customvalidator "baneks.com/internal/utils/validator"
	aiClient "baneks.com/pkg/ai"
	"github.com/labstack/echo/v5"
)

type handlerRequest struct {
	Text string `json:"text" validate:"required,max=100"`
}

type ChatWithAiHandler struct {
	aiClient *aiClient.AiClient
}

func NewChatWithAiHandler(aiClient *aiClient.AiClient) *ChatWithAiHandler {
	return &ChatWithAiHandler{
		aiClient: aiClient,
	}
}

func (h *ChatWithAiHandler) ChatWithAi(c *echo.Context) error {
	params := &handlerRequest{}

	if err := c.Bind(params); err != nil {
		return customerrors.NewAppBindError(err)
	}

	if err := customvalidator.ValidateRequest(c.Validate, params); err != nil {
		return customerrors.NewAppHTTPError(http.StatusForbidden, err, nil)
	}

	resp, err := h.aiClient.Chat(c.Request().Context(), aiClient.ChatRequest{
		MaxResponseTokens: 200,
		Messages: []aiClient.ChatMessage{
			{
				Role:    "system",
				Content: "Всегда пиши на русском языке и plain text. Не используй markdown разметку. Очень коротко, 2-3 предложений.",
			},
			{
				Role:    "user",
				Content: params.Text,
			},
		},
	})
	if err != nil {
		if errors.Is(err, aiClient.ErrRateLimitExceeded) {
			return customerrors.NewAppHTTPError(http.StatusTooManyRequests, "To many requests", err)
		}

		return customerrors.NewAppHTTPError(http.StatusInternalServerError, "Something went wrong", err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"text": resp,
	})
}
