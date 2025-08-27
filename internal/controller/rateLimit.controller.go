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
		slog.Warn(
			"Failed to parse rate limit header",
			"error", err,
			"header", resp.Header.Get("X-Discogs-Ratelimit"),
		)
		rateLimit = 0 // Default value on parse error
	}

	rateLimitUsed, err := strconv.Atoi(resp.Header.Get("X-Discogs-Ratelimit-Used"))
	if err != nil {
		slog.Warn(
			"Failed to parse rate limit used header",
			"error", err,
			"header", resp.Header.Get("X-Discogs-Ratelimit-Used"),
		)
		rateLimitUsed = 0 // Default value on parse error
	}

	rateLimitRemaining, err := strconv.Atoi(resp.Header.Get("X-Discogs-Ratelimit-Remaining"))
	if err != nil {
		slog.Warn(
			"Failed to parse rate limit remaining header",
			"error", err,
			"header", resp.Header.Get("X-Discogs-Ratelimit-Remaining"),
		)
		rateLimitRemaining = 0 // Default value on parse error
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	oldLimits := RateLimit{Limit: r.Limit, Used: r.Used, Remaining: r.Remaining}

	r.Limit = rateLimit
	r.Used = rateLimitUsed
	r.Remaining = rateLimitRemaining

	slog.Debug("Updated rate limits", 
		"oldLimit", oldLimits.Limit,
		"oldUsed", oldLimits.Used,
		"oldRemaining", oldLimits.Remaining,
		"newLimit", r.Limit,
		"newUsed", r.Used,
		"newRemaining", r.Remaining)

	if r.Remaining <= 10 {
		slog.Warn("Rate limit threshold reached", 
			"remaining", r.Remaining,
			"used", r.Used,
			"limit", r.Limit,
			"shouldThrottle", true)
	}
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
