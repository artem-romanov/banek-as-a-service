package ai

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Limits for GPT OSS 120b
// https://inference-docs.cerebras.ai/support/rate-limits
const (
	LimitRequestsMinute = 30
	LimitRequestsHour   = 900
	LimitRequestsDay    = 14_400

	// LimitTokensMinute = 64_000
	LimitTokensMinute = 30_000
	LimitTokensHour   = 1_000_000
	LimitTokensDay    = 1_000_000
)

// RateLimits предоставляет информацию о лимитах доступных на данный момент
// Так же имеет информацию о том когда в последний раз был выполнен запрос(ы)
type RateLimits struct {
	RequestsMinute int
	RequestsHour   int
	RequestsDay    int

	TokensMinute int
	TokensHour   int
	TokensDay    int

	MinuteStartedAt time.Time
	HourStartedAt   time.Time
	DayStartedAt    time.Time

	mu sync.Mutex
}

func DefaultRateLimits() *RateLimits {
	rl := &RateLimits{}
	rl.resetIfNeeded(time.Now())
	return rl
}

// TryReserve does two things in transaction
// - Check can a new request be done
// - If it can be done - it reserves some tokens for the request
// Function helps to avoid concurrency problems:
// e.g. 2 parallel requests are allowed to do a request and using same limits
//
// Warn: function highly depends on counted tokens.
// If you counted them wrong - you might be allowed to do a request and get 429!
func (l *RateLimits) TryReserve(tokens int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	l.resetIfNeeded(now)

	if l.RequestsMinute+1 > LimitRequestsMinute ||
		l.RequestsHour+1 > LimitRequestsHour ||
		l.RequestsDay+1 > LimitRequestsDay ||
		l.TokensMinute+tokens > LimitTokensMinute ||
		l.TokensHour+tokens > LimitTokensHour ||
		l.TokensDay+tokens > LimitTokensDay {
		return false
	}

	l.RequestsMinute++
	l.RequestsHour++
	l.RequestsDay++

	l.TokensMinute += tokens
	l.TokensHour += tokens
	l.TokensDay += tokens

	return true
}

// resetIfNeeded checks and resets limits in Allow()
// Warn: should run inside Lock!
func (l *RateLimits) resetIfNeeded(now time.Time) {
	if now.Sub(l.MinuteStartedAt) >= time.Minute {
		l.RequestsMinute = 0
		l.TokensMinute = 0
		l.MinuteStartedAt = now
	}

	if now.Sub(l.HourStartedAt) >= time.Hour {
		l.RequestsHour = 0
		l.TokensHour = 0
		l.HourStartedAt = now
	}

	if now.Sub(l.DayStartedAt) >= 24*time.Hour {
		l.RequestsDay = 0
		l.TokensDay = 0
		l.DayStartedAt = now
	}
}

func (l *RateLimits) UpdateFromHeaders(h http.Header) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	// Минутное окно
	if rem, ok := headerInt(h, "x-ratelimit-remaining-requests-minute"); ok {
		used := LimitRequestsMinute - rem

		// Updating only if current record was created LATER than we have in rate limit.
		// Here and in further it saves us from concurrency hell (e.g. earliest request was the slowest)
		if used > l.RequestsMinute || now.Sub(l.MinuteStartedAt) >= time.Minute {
			l.RequestsMinute = used
		}
	}
	if remTokens, ok := headerInt(h, "x-ratelimit-remaining-tokens-minute"); ok {
		if remTokens < l.TokensMinute || now.Sub(l.MinuteStartedAt) >= time.Minute {
			l.TokensMinute = LimitTokensMinute - remTokens
			l.MinuteStartedAt = now
		}
	}

	// Часовое окно
	if rem, ok := headerInt(h, "x-ratelimit-remaining-requests-hour"); ok {
		used := LimitRequestsHour - rem
		if used > l.RequestsHour || now.Sub(l.HourStartedAt) >= time.Hour {
			l.RequestsHour = used
		}
	}
	if remTokens, ok := headerInt(h, "x-ratelimit-remaining-tokens-hour"); ok {
		if remTokens < l.TokensHour || now.Sub(l.HourStartedAt) >= time.Hour {
			l.TokensHour = LimitTokensHour - remTokens
			l.HourStartedAt = now
		}
	}

	// Дневное окно
	if rem, ok := headerInt(h, "x-ratelimit-remaining-requests-day"); ok {
		used := LimitRequestsDay - rem
		if used > l.RequestsDay || now.Sub(l.DayStartedAt) >= 24*time.Hour {
			l.RequestsDay = used
		}
	}
	if remTokens, ok := headerInt(h, "x-ratelimit-remaining-tokens-day"); ok {
		if remTokens < l.TokensDay || now.Sub(l.DayStartedAt) >= 24*time.Hour {
			l.TokensDay = LimitTokensDay - remTokens
			l.DayStartedAt = now
		}
	}
}

func headerInt(h http.Header, key string) (int, bool) {
	v := h.Get(key)
	if v == "" {
		return 0, false
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return i, true
}
