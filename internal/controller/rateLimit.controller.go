package controller

import (
	"log/slog"
	"net/http"
	"strconv"
	"sync"
)

type RateLimit struct {
	mutex     sync.RWMutex
	Limit     int `json:"limit"`
	Used      int `json:"used"`
	Remaining int `json:"remaining"`
}

func (r *RateLimit) UpdateLimits(resp *http.Response) {
	rateLimit, err := strconv.Atoi(resp.Header.Get("X-Discogs-Ratelimit"))
	if err != nil {
		slog.Error(
			"Failed to parse rate limit header",
			"error", err,
			"header", resp.Header.Get("X-Discogs-Ratelimit"),
		)
	}

	rateLimitUsed, err := strconv.Atoi(resp.Header.Get("X-Discogs-Ratelimit-Used"))
	if err != nil {
		slog.Error(
			"Failed to parse rate limit used header",
			"error", err,
			"header", resp.Header.Get("X-Discogs-Ratelimit-Used"),
		)
	}

	rateLimitRemaining, err := strconv.Atoi(resp.Header.Get("X-Discogs-Ratelimit-Remaining"))
	if err != nil {
		slog.Error(
			"Failed to parse rate limit remaining header",
			"error", err,
			"header", resp.Header.Get("X-Discogs-Ratelimit-Remaining"),
		)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.Limit = rateLimit
	r.Used = rateLimitUsed
	r.Remaining = rateLimitRemaining
}

func (r *RateLimit) GetCurrent() RateLimit {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return RateLimit{
		Limit:     r.Limit,
		Used:      r.Used,
		Remaining: r.Remaining,
	}
}

func (r *RateLimit) ShouldThrottle() bool {
	return r.GetCurrent().Remaining <= 10
}
