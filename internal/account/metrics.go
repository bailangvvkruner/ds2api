package account

import "time"

type usageEvent struct {
	At     time.Time
	Tokens int
}

type AccountMetrics struct {
	TodayRequests     int   `json:"today_requests"`
	TodayInputTokens  int64 `json:"today_input_tokens"`
	TodayOutputTokens int64 `json:"today_output_tokens"`
	TotalInputTokens  int64 `json:"total_input_tokens"`
	TotalOutputTokens int64 `json:"total_output_tokens"`
	RPM               int   `json:"rpm"`
	TPM               int   `json:"tpm"`
	AverageResponseMs int64 `json:"average_response_ms"`
	AverageTimeMs     int64 `json:"average_time_ms"`
}

type accountMetricsState struct {
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

func (p *Pool) SetPaused(accountID string, paused bool) bool {
	if accountID == "" {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.store.FindAccount(accountID); !ok {
		return false
	}
	if p.paused == nil {
		p.paused = map[string]bool{}
	}
	if paused {
		p.paused[accountID] = true
	} else {
		delete(p.paused, accountID)
	}
	p.notifyWaiterLocked()
	return true
}

func (p *Pool) RecordUsage(accountID string, inputTokens, outputTokens int, elapsed time.Duration) {
	if accountID == "" {
		return
	}
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
	day := now.Local().Format("2006-01-02")
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.metrics == nil {
		p.metrics = map[string]*accountMetricsState{}
	}
	m := p.metrics[accountID]
	if m == nil {
		m = &accountMetricsState{day: day}
		p.metrics[accountID] = m
	}
	if m.day != day {
		m.day = day
		m.todayRequests = 0
		m.todayInputTokens = 0
		m.todayOutputTokens = 0
	}
	m.todayRequests++
	m.todayInputTokens += int64(inputTokens)
	m.todayOutputTokens += int64(outputTokens)
	m.totalInputTokens += int64(inputTokens)
	m.totalOutputTokens += int64(outputTokens)
	m.completedRequests++
	elapsedMs := elapsed.Milliseconds()
	m.totalResponseMs += elapsedMs
	m.totalTimeMs += elapsedMs
	m.events = append(m.events, usageEvent{At: now, Tokens: inputTokens + outputTokens})
	m.events = trimUsageEvents(now, m.events)
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

func (p *Pool) metricsSnapshotLocked(accountID string, now time.Time) AccountMetrics {
	if p.metrics == nil {
		return AccountMetrics{}
	}
	m := p.metrics[accountID]
	if m == nil {
		return AccountMetrics{}
	}
	day := now.Local().Format("2006-01-02")
	todayRequests := m.todayRequests
	todayInputTokens := m.todayInputTokens
	todayOutputTokens := m.todayOutputTokens
	if m.day != day {
		todayRequests = 0
		todayInputTokens = 0
		todayOutputTokens = 0
	}
	events := trimUsageEvents(now, m.events)
	m.events = events
	tpm := 0
	for _, event := range events {
		tpm += event.Tokens
	}
	avgResponseMs := int64(0)
	avgTimeMs := int64(0)
	if m.completedRequests > 0 {
		avgResponseMs = m.totalResponseMs / m.completedRequests
		avgTimeMs = m.totalTimeMs / m.completedRequests
	}
	return AccountMetrics{
		TodayRequests:     todayRequests,
		TodayInputTokens:  int64(todayInputTokens),
		TodayOutputTokens: int64(todayOutputTokens),
		TotalInputTokens:  m.totalInputTokens,
		TotalOutputTokens: m.totalOutputTokens,
		RPM:               len(events),
		TPM:               tpm,
		AverageResponseMs: avgResponseMs,
		AverageTimeMs:     avgTimeMs,
	}
}
