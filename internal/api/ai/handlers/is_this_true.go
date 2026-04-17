package handlers

import (
	"errors"
	"math/rand"
	"net/http"
	"strings"

	customerrors "baneks.com/internal/custom_errors"
	customvalidator "baneks.com/internal/utils/validator"
	aiClient "baneks.com/pkg/ai"
	"github.com/labstack/echo/v5"
)

var prompts = []string{
	"Дай ответ на вопрос \"это правда\" следующего сообщения так будто ты знаток всего. Не более 3-4 предложений.",
}

type isThisTrueRequest struct {
	Question string `json:"question" validate:"required,max=3000"`
}

// IsThisTrueHanlder very specific ai hanlder which is needed for grok-like requests.
// Although we could try to pollute "chat with ai" handler, I decided to separate them.
type IsThisTrueHandler struct {
	aiClient *aiClient.AiClient
}

func NewIsThisTrueHandler(aiClient *aiClient.AiClient) *IsThisTrueHandler {
	return &IsThisTrueHandler{
		aiClient: aiClient,
	}
}

func (h *IsThisTrueHandler) IsThisTrueHandler(c *echo.Context) error {
	params := new(isThisTrueRequest)

	if err := c.Bind(params); err != nil {
		return customerrors.NewAppBindError(err)
	}

	if err := customvalidator.ValidateRequest(c.Validate, params); err != nil {
		return customerrors.NewAppHTTPError(http.StatusBadRequest, err, nil)
	}

	sanitizedQuestion := strings.TrimSpace(params.Question)
	if len(sanitizedQuestion) == 0 {
		return customerrors.NewAppHTTPError(http.StatusBadRequest, "question is missing", nil)
	}

	resp, err := h.aiClient.Chat(c.Request().Context(), aiClient.ChatRequest{
		MaxResponseTokens: 250,
		Messages: []aiClient.ChatMessage{
			{
				Role:    "system",
				Content: getSystemPrompt(),
			},
			{
				Role:    "user",
				Content: sanitizedQuestion,
			},
		},
	})
	if err != nil {
		if errors.Is(err, aiClient.ErrRateLimitExceeded) {
			return customerrors.NewAppHTTPError(http.StatusTooManyRequests, "Too many requests", err)
		}

		return customerrors.NewAppHTTPError(http.StatusInternalServerError, "Something went wrong", err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"text": resp,
	})
}

func getSystemPrompt() string {
	index := rand.Intn(len(prompts))
	return prompts[index]
}
