package account

import (
	"sync"
	"time"
)

type Statistics struct {
	mu                sync.Mutex
	TodayRequests     int64
	TodayInputTokens  int64
	TodayOutputTokens int64
	TotalInputTokens  int64
	TotalOutputTokens int64
	TodayStartTime    time.Time
	RecentRequests    []RequestRecord
	RPMWindow         []time.Time
	TPMWindow         []int64
}

type RequestRecord struct {
	Timestamp    time.Time
	InputTokens  int64
	OutputTokens int64
	Duration     time.Duration
}

var globalStats = &Statistics{
	TodayStartTime: getTodayStart(),
	RecentRequests: make([]RequestRecord, 0, 1000),
	RPMWindow:      make([]time.Time, 0, 100),
	TPMWindow:      make([]int64, 0, 100),
}

func getTodayStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func (s *Statistics) RecordRequest(inputTokens, outputTokens int64, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if now.Sub(s.TodayStartTime) >= 24*time.Hour {
		s.TodayRequests = 0
		s.TodayInputTokens = 0
		s.TodayOutputTokens = 0
		s.TodayStartTime = getTodayStart()
	}

	s.TodayRequests++
	s.TodayInputTokens += inputTokens
	s.TodayOutputTokens += outputTokens
	s.TotalInputTokens += inputTokens
	s.TotalOutputTokens += outputTokens

	s.RPMWindow = append(s.RPMWindow, now)
	s.TPMWindow = append(s.TPMWindow, inputTokens+outputTokens)

	cutoff := now.Add(-time.Minute)
	for len(s.RPMWindow) > 0 && s.RPMWindow[0].Before(cutoff) {
		s.RPMWindow = s.RPMWindow[1:]
	}

	cutoff = now.Add(-time.Minute)
	for len(s.TPMWindow) > 0 && len(s.RPMWindow) > 0 && s.RPMWindow[0].Before(cutoff) {
		s.TPMWindow = s.TPMWindow[1:]
		s.RPMWindow = s.RPMWindow[1:]
	}

	if len(s.RecentRequests) >= 1000 {
		s.RecentRequests = s.RecentRequests[1:]
	}
	s.RecentRequests = append(s.RecentRequests, RequestRecord{
		Timestamp:    now,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		Duration:     duration,
	})
}

func (s *Statistics) GetRPM() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return int64(len(s.RPMWindow))
}

func (s *Statistics) GetTPM() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	var total int64
	for _, t := range s.TPMWindow {
		total += t
	}
	return total
}

func (s *Statistics) GetAverageResponseTime() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.RecentRequests) == 0 {
		return 0
	}
	var total int64
	for _, r := range s.RecentRequests {
		total += r.Duration.Milliseconds()
	}
	return float64(total) / float64(len(s.RecentRequests)) / 1000.0
}

func RecordRequest(inputTokens, outputTokens int64, duration time.Duration) {
	globalStats.RecordRequest(inputTokens, outputTokens, duration)
}

func GetStatistics() map[string]any {
	globalStats.mu.Lock()
	defer globalStats.mu.Unlock()

	now := time.Now()
	if now.Sub(globalStats.TodayStartTime) >= 24*time.Hour {
		return map[string]any{
			"today_requests":      int64(0),
			"today_input_tokens":  int64(0),
			"today_output_tokens": int64(0),
			"total_input_tokens":  globalStats.TotalInputTokens,
			"total_output_tokens": globalStats.TotalOutputTokens,
			"rpm":                 int64(0),
			"tpm":                 int64(0),
			"avg_response_time":   0.0,
		}
	}

	var avgResponseTime float64
	if len(globalStats.RecentRequests) > 0 {
		var total int64
		for _, r := range globalStats.RecentRequests {
			total += r.Duration.Milliseconds()
		}
		avgResponseTime = float64(total) / float64(len(globalStats.RecentRequests)) / 1000.0
	}

	return map[string]any{
		"today_requests":      globalStats.TodayRequests,
		"today_input_tokens":  globalStats.TodayInputTokens,
		"today_output_tokens": globalStats.TodayOutputTokens,
		"total_input_tokens":  globalStats.TotalInputTokens,
		"total_output_tokens": globalStats.TotalOutputTokens,
		"rpm":                 int64(len(globalStats.RPMWindow)),
		"tpm":                 globalStats.GetTPM(),
		"avg_response_time":   avgResponseTime,
	}
}
