package ai

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type RateLimitStorage interface {
	Save(*RateLimits) error
	Load() (*RateLimits, error)
}

type jsonRateLimitStorage struct {
	Path string
	mu   sync.Mutex
}

func NewJsonRateLimitStorage() *jsonRateLimitStorage {
	return &jsonRateLimitStorage{
		Path: "limits.json",
	}

}

func (s *jsonRateLimitStorage) Save(limit *RateLimits) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Create(s.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// сохраняем только данные, без mutex
	data := struct {
		RequestsMinute  int       `json:"requests_minute"`
		RequestsHour    int       `json:"requests_hour"`
		RequestsDay     int       `json:"requests_day"`
		TokensMinute    int       `json:"tokens_minute"`
		TokensHour      int       `json:"tokens_hour"`
		TokensDay       int       `json:"tokens_day"`
		MinuteStartedAt time.Time `json:"minute_started_at"`
		HourStartedAt   time.Time `json:"hour_started_at"`
		DayStartedAt    time.Time `json:"day_started_at"`
	}{
		RequestsMinute:  limit.RequestsMinute,
		RequestsHour:    limit.RequestsHour,
		RequestsDay:     limit.RequestsDay,
		TokensMinute:    limit.TokensMinute,
		TokensHour:      limit.TokensHour,
		TokensDay:       limit.TokensDay,
		MinuteStartedAt: limit.MinuteStartedAt,
		HourStartedAt:   limit.HourStartedAt,
		DayStartedAt:    limit.DayStartedAt,
	}

	return json.NewEncoder(file).Encode(data)
}

func (s *jsonRateLimitStorage) Load() (*RateLimits, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var limit RateLimits
	file, err := os.Open(s.Path)
	if err != nil {
		if os.IsNotExist(err) {
			// файл не существует → возвращаем пустые лимиты
			return &RateLimits{}, nil
		}
		return &RateLimits{}, err
	}
	defer file.Close()

	data := struct {
		RequestsMinute  int       `json:"requests_minute"`
		RequestsHour    int       `json:"requests_hour"`
		RequestsDay     int       `json:"requests_day"`
		TokensMinute    int       `json:"tokens_minute"`
		TokensHour      int       `json:"tokens_hour"`
		TokensDay       int       `json:"tokens_day"`
		MinuteStartedAt time.Time `json:"minute_started_at"`
		HourStartedAt   time.Time `json:"hour_started_at"`
		DayStartedAt    time.Time `json:"day_started_at"`
	}{}

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	limit.RequestsMinute = data.RequestsMinute
	limit.RequestsHour = data.RequestsHour
	limit.RequestsDay = data.RequestsDay
	limit.TokensMinute = data.TokensMinute
	limit.TokensHour = data.TokensHour
	limit.TokensDay = data.TokensDay
	limit.MinuteStartedAt = data.MinuteStartedAt
	limit.HourStartedAt = data.HourStartedAt
	limit.DayStartedAt = data.DayStartedAt

	return &limit, nil
}
