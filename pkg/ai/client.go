package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AiClient struct {
	client   *http.Client
	sem      chan struct{}
	provider Provider
}

// NewAiClient creates new client with selected Provider
// It panics if provider is nil
func NewAiClient(
	client *http.Client,
	provider Provider,
	maxConcurrentLimit int,
) *AiClient {
	semLimit := 1
	if maxConcurrentLimit > 0 {
		semLimit = maxConcurrentLimit
	}
	if provider == nil {
		panic("provider is nil")
	}
	return &AiClient{
		client:   client,
		sem:      make(chan struct{}, semLimit),
		provider: provider,
	}
}

type chatRequestMessagesDTO struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type chatRequestDTO struct {
	Model           string                   `json:"model"`
	Stream          bool                     `json:"stream"`
	MaxTokens       int                      `json:"max_tokens"`
	Temperature     float64                  `json:"temperature"`
	ReasoningEffort string                   `json:"reasoning_effort"`
	Messages        []chatRequestMessagesDTO `json:"messages"`
}

type chatResponseDTO struct {
	Choices []struct {
		Message struct {
			Content   string `json:"content"`
			Reasoning string `json:"reasoning"`
		} `json:"message"`
	} `json:"choices"`
}

func (r chatResponseDTO) FirstChoice() (content, reasoning string) {
	if len(r.Choices) == 0 {
		return "", ""
	}
	msg := r.Choices[0].Message
	return msg.Content, msg.Reasoning
}

func (c *AiClient) Chat(
	ctx context.Context,
	request ChatRequest,
) (string, error) {
	select {
		case c.sem <- struct{}{}:
			defer func() { <-c.sem}()
		case <-ctx.Done():
			return "", ctx.Err()
	}

	dto := chatRequestToDTO(request, c.provider.Model())
	requestBody, err := json.Marshal(dto)
	if err != nil {
		return "", fmt.Errorf("request body marshall error: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.provider.BaseUrl()+"/v1/chat/completions",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return "", fmt.Errorf("create http request error: %w", err)
	}
	c.withAuthorization(httpReq)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("chat request error: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cant parse request body: %w, status code = %v", err, resp.StatusCode)
	}

	if resp.StatusCode == 429 {
		return "", ErrRateLimitExceeded
	} else if resp.StatusCode >= 400 {
		return "", fmt.Errorf("chat request status is %v, message=%v", resp.StatusCode, string(responseBody))
	}

	var respDto chatResponseDTO
	if err := json.Unmarshal(responseBody, &respDto); err != nil {
		return "", fmt.Errorf("cant unmarshall chat completion result: %w", err)
	}

	// we dont need reasoning for now
	content, _ := respDto.FirstChoice()
	return content, nil
}

func (c *AiClient) withAuthorization(request *http.Request) {
	request.Header.Add("Authorization", "Bearer "+c.provider.Token())
}

func chatRequestToDTO(req ChatRequest, model string) chatRequestDTO {
	messages := make([]chatRequestMessagesDTO, len(req.Messages))

	for i, msg := range req.Messages {
		messages[i].Role = msg.Role
		messages[i].Content = msg.Content
	}

	return chatRequestDTO{
		Model:           model,
		Stream:          false,
		MaxTokens:       req.MaxResponseTokens,
		Temperature:     1,
		ReasoningEffort: "low",
		Messages:        messages,
	}
}
