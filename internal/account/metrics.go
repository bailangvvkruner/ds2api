package account

import "time"

type usageEvent struct {
	At     time.Time
	Tokens int
}

type Metrics struct {
	day               string
	todayRequests     int
	todayInputTokens  int64
	todayOutputTokens int64
	totalInputTokens  int64
	totalOutputTokens int64
	totalResponseMs   int64
	totalTimeMs       int64
	completedRequests int64
	events            []usageEvent
}

func NewMetrics() *Metrics {
	return &Metrics{day: time.Now().Local().Format("2006-01-02")}
}

func (p *Pool) RecordUsage(_ string, inputTokens, outputTokens int, elapsed time.Duration) {
	if inputTokens < 0 {
		inputTokens = 0
	}
	if outputTokens < 0 {
		outputTokens = 0
	}
	if elapsed < 0 {
		elapsed = 0
	}
	now := time.Now()
	p.mu.Lock()
	defer p.mu.Unlock()
	p.recordUsageLocked(now, inputTokens, outputTokens, elapsed)
}

func (p *Pool) recordUsageLocked(now time.Time, inputTokens, outputTokens int, elapsed time.Duration) {
	if p.metrics == nil {
		p.metrics = NewMetrics()
	}
	day := now.Local().Format("2006-01-02")
	if p.metrics.day != day {
		p.metrics.day = day
		p.metrics.todayRequests = 0
		p.metrics.todayInputTokens = 0
		p.metrics.todayOutputTokens = 0
	}
	p.metrics.todayRequests++
	p.metrics.todayInputTokens += int64(inputTokens)
	p.metrics.todayOutputTokens += int64(outputTokens)
	p.metrics.totalInputTokens += int64(inputTokens)
	p.metrics.totalOutputTokens += int64(outputTokens)
	p.metrics.completedRequests++
	elapsedMs := elapsed.Milliseconds()
	p.metrics.totalResponseMs += elapsedMs
	p.metrics.totalTimeMs += elapsedMs
	p.metrics.events = append(p.metrics.events, usageEvent{At: now, Tokens: inputTokens + outputTokens})
	p.metrics.events = trimUsageEvents(now, p.metrics.events)
}

func trimUsageEvents(now time.Time, events []usageEvent) []usageEvent {
	cutoff := now.Add(-time.Minute)
	idx := 0
	for idx < len(events) && events[idx].At.Before(cutoff) {
		idx++
	}
	if idx == 0 {
		return events
	}
	trimmed := events[idx:]
	out := make([]usageEvent, len(trimmed))
	copy(out, trimmed)
	return out
}

func (p *Pool) metricsStatusLocked(now time.Time) map[string]any {
	if p.metrics == nil {
		p.metrics = NewMetrics()
	}
	day := now.Local().Format("2006-01-02")
	todayRequests := p.metrics.todayRequests
	todayInputTokens := p.metrics.todayInputTokens
	todayOutputTokens := p.metrics.todayOutputTokens
	if p.metrics.day != day {
		todayRequests = 0
		todayInputTokens = 0
		todayOutputTokens = 0
	}
	events := trimUsageEvents(now, p.metrics.events)
	p.metrics.events = events
	tpm := 0
	for _, event := range events {
		tpm += event.Tokens
	}
	averageResponseMs := int64(0)
	averageTimeMs := int64(0)
	if p.metrics.completedRequests > 0 {
		averageResponseMs = p.metrics.totalResponseMs / p.metrics.completedRequests
		averageTimeMs = p.metrics.totalTimeMs / p.metrics.completedRequests
	}
	return map[string]any{
		"today_requests":      todayRequests,
		"today_input_tokens":  todayInputTokens,
		"today_output_tokens": todayOutputTokens,
		"total_input_tokens":  p.metrics.totalInputTokens,
		"total_output_tokens": p.metrics.totalOutputTokens,
		"rpm":                 len(events),
		"tpm":                 tpm,
		"average_response_ms": averageResponseMs,
		"average_time_ms":     averageTimeMs,
	}
}
