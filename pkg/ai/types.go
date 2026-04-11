package ai

import "errors"

// ERRORS
var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// Domain package models
type ChatRequest struct {
	MaxResponseTokens int
	Messages          []ChatMessage
}

type ChatMessage struct {
	Role    string // system, user
	Content string // prompt
}
