package account

import (
	"sync"
	"time"

	"ds2api/internal/config"
)

type Stats struct {
	mu                sync.RWMutex
	todayRequests     int64
	todayInputTokens  int64
	todayOutputTokens int64
	totalInputTokens  int64
	totalOutputTokens int64
	avgResponseTime   float64
	responseCount     int64
	requestTimes      []time.Time
	tokenCounts       []int64
	windowDuration    time.Duration
	lastResetDay      int
}

func NewStats() *Stats {
	s := &Stats{
		windowDuration: time.Minute,
		lastResetDay:   time.Now().Day(),
	}
	return s
}

func (s *Stats) RecordRequest(inputTokens, outputTokens int64, responseTime float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	config.Logger.Debug("[stats] RecordRequest called", "inputTokens", inputTokens, "outputTokens", outputTokens, "responseTime", responseTime)

	now := time.Now()
	if now.Day() != s.lastResetDay {
		s.todayRequests = 0
		s.todayInputTokens = 0
		s.todayOutputTokens = 0
		s.lastResetDay = now.Day()
	}

	s.todayRequests++
	s.todayInputTokens += inputTokens
	s.todayOutputTokens += outputTokens
	s.totalInputTokens += inputTokens
	s.totalOutputTokens += outputTokens

	s.requestTimes = append(s.requestTimes, now)
	s.tokenCounts = append(s.tokenCounts, inputTokens+outputTokens)

	if responseTime > 0 {
		totalTime := s.avgResponseTime*float64(s.responseCount) + responseTime
		s.responseCount++
		s.avgResponseTime = totalTime / float64(s.responseCount)
	}
}

func (s *Stats) GetStatus() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	rpm := s.calculateRPM(now)
	tpm := s.calculateTPM(now)

	return map[string]any{
		"today_requests":      s.todayRequests,
		"today_input_tokens":  s.todayInputTokens,
		"today_output_tokens": s.todayOutputTokens,
		"total_input_tokens":  s.totalInputTokens,
		"total_output_tokens": s.totalOutputTokens,
		"rpm":                 rpm,
		"tpm":                 tpm,
		"avg_response_time":   s.avgResponseTime,
	}
}

func (s *Stats) calculateRPM(now time.Time) int64 {
	cutoff := now.Add(-s.windowDuration)
	count := int64(0)
	for _, t := range s.requestTimes {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}

func (s *Stats) calculateTPM(now time.Time) int64 {
	cutoff := now.Add(-s.windowDuration)
	total := int64(0)
	for i, t := range s.requestTimes {
		if t.After(cutoff) && i < len(s.tokenCounts) {
			total += s.tokenCounts[i]
		}
	}
	return total
}
