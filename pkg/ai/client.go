package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"unicode/utf8"
)

const modelName = "llama3.1-8b"
const apiUrl = "https://api.cerebras.ai/v1/chat/completions"

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

type ChatRequest struct {
	MaxResponseTokens int
	Messages          []ChatMessage
}

type ChatMessage struct {
	Role    string // system, user
	Content string // prompt
}

func (c ChatRequest) EstimatePrompt() int {
	total := 0
	for _, msg := range c.Messages {
		total += utf8.RuneCountInString(msg.Content)
	}
	// conservative estimation, 1 token = 3 symbols
	// TODO: think about GPT OSS tokenizer, we need o200k_harmony
	// https://cdn.openai.com/pdf/419b6906-9da6-406c-a19d-1bb078ac7637/oai_gpt-oss_model_card.pdf
	return total/3 + c.MaxResponseTokens
}

type AiClient struct {
	client       *http.Client
	limits       *RateLimits
	token        string
	sem          chan struct{}
	limitStorage RateLimitStorage
}

func NewAiClient(
	client *http.Client,
	token string,
	maxConcurrentLimit int,
	storage RateLimitStorage,
) (*AiClient, error) {
	limits := DefaultRateLimits()
	if storage != nil {
		loadedLimits, err := storage.Load()
		if err != nil {
			return nil, fmt.Errorf("rate limit load error: %w", err)
		} else {
			limits = loadedLimits
		}
	} else {
		storage = NewJsonRateLimitStorage()
	}

	return &AiClient{
		client:       client,
		limits:       limits,
		token:        token,
		sem:          make(chan struct{}, maxConcurrentLimit),
		limitStorage: storage,
	}, nil
}

type chatRequestMessagesDTO struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type chatRequestDTO struct {
	Model       string  `json:"model"`
	Stream      bool    `json:"stream"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	// ReasoningEffort string                   `json:"reasoning_effort"`
	Messages []chatRequestMessagesDTO `json:"messages"`
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
	tokensEstimate := request.EstimatePrompt()

	if !c.limits.TryReserve(tokensEstimate) {
		return "", ErrRateLimitExceeded
	}

	c.sem <- struct{}{}
	defer func() { <-c.sem }()

	dto := chatRequestToDTO(request)
	requestBody, err := json.Marshal(dto)
	if err != nil {
		return "", fmt.Errorf("request body marshall error: %w", err)
	}

	// TODO: I don't like bytes.NewReader(body).
	// Think about how to decrease allocations.
	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiUrl, bytes.NewReader(requestBody))
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

	c.limits.UpdateFromHeaders(resp.Header)
	if err := c.limitStorage.Save(c.limits); err != nil {
		// TODO: think about better way to work around this
		fmt.Println("couldnt save limits on disk", err)
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
	request.Header.Add("Authorization", "Bearer "+c.token)
}

func chatRequestToDTO(req ChatRequest) chatRequestDTO {
	messages := make([]chatRequestMessagesDTO, len(req.Messages))

	for i, msg := range req.Messages {
		messages[i].Role = msg.Role
		messages[i].Content = msg.Content
	}

	return chatRequestDTO{
		Model:       modelName,
		Stream:      false,
		MaxTokens:   req.MaxResponseTokens,
		Temperature: 1,
		TopP:        1,
		// ReasoningEffort: "low",
		Messages: messages,
	}
}
