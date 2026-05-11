package account

import (
	"sort"
	"sync"
	"time"

	"ds2api/internal/config"
)

type Pool struct {
	store                  *config.Store
	mu                     sync.Mutex
	queue                  []string
	inUse                  map[string]int
	paused                 map[string]bool
	metrics                map[string]*accountMetricsState
	waiters                []chan struct{}
	maxInflightPerAccount  int
	recommendedConcurrency int
	maxQueueSize           int
	globalMaxInflight      int
}

func NewPool(store *config.Store) *Pool {
	maxPer := 2
	if store != nil {
		maxPer = store.RuntimeAccountMaxInflight()
	}
	p := &Pool{
		store:                 store,
		inUse:                 map[string]int{},
		paused:                map[string]bool{},
		metrics:               map[string]*accountMetricsState{},
		maxInflightPerAccount: maxPer,
	}
	p.Reset()
	return p
}

func (p *Pool) Reset() {
	accounts := p.store.Accounts()
	sort.SliceStable(accounts, func(i, j int) bool {
		iHas := accounts[i].Token != ""
		jHas := accounts[j].Token != ""
		if iHas == jHas {
			return i < j
		}
		return iHas
	})
	ids := make([]string, 0, len(accounts))
	for _, a := range accounts {
		id := a.Identifier()
		if id != "" {
			ids = append(ids, id)
		}
	}
	if p.store != nil {
		p.maxInflightPerAccount = p.store.RuntimeAccountMaxInflight()
	} else {
		p.maxInflightPerAccount = maxInflightFromEnv()
	}
	recommended := defaultRecommendedConcurrency(len(ids), p.maxInflightPerAccount)
	queueLimit := maxQueueFromEnv(recommended)
	globalLimit := recommended
	if p.store != nil {
		queueLimit = p.store.RuntimeAccountMaxQueue(recommended)
		globalLimit = p.store.RuntimeGlobalMaxInflight(recommended)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.drainWaitersLocked()
	p.queue = ids
	p.inUse = map[string]int{}
	if p.paused == nil {
		p.paused = map[string]bool{}
	}
	for id := range p.paused {
		found := false
		for _, queueID := range ids {
			if id == queueID {
				found = true
				break
			}
		}
		if !found {
			delete(p.paused, id)
		}
	}
	p.recommendedConcurrency = recommended
	p.maxQueueSize = queueLimit
	p.globalMaxInflight = globalLimit
	config.Logger.Info(
		"[init_account_queue] initialized",
		"total", len(ids),
		"max_inflight_per_account", p.maxInflightPerAccount,
		"global_max_inflight", p.globalMaxInflight,
		"recommended_concurrency", p.recommendedConcurrency,
		"max_queue_size", p.maxQueueSize,
	)
}

func (p *Pool) Release(accountID string) {
	if accountID == "" {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	count := p.inUse[accountID]
	if count <= 0 {
		return
	}
	if count == 1 {
		delete(p.inUse, accountID)
		p.notifyWaiterLocked()
		return
	}
	p.inUse[accountID] = count - 1
	p.notifyWaiterLocked()
}

func (p *Pool) Status() map[string]any {
	p.mu.Lock()
	defer p.mu.Unlock()
	available := make([]string, 0, len(p.queue))
	inUseAccounts := make([]string, 0, len(p.inUse))
	pausedAccounts := make([]string, 0, len(p.paused))
	accountStats := make(map[string]any, len(p.queue))
	inUseSlots := 0
	now := time.Now()
	for _, id := range p.queue {
		metrics := p.metricsSnapshotLocked(id, now)
		accountStats[id] = map[string]any{
			"paused":              p.paused[id],
			"in_use":              p.inUse[id],
			"today_requests":      metrics.TodayRequests,
			"today_input_tokens":  metrics.TodayInputTokens,
			"today_output_tokens": metrics.TodayOutputTokens,
			"total_input_tokens":  metrics.TotalInputTokens,
			"total_output_tokens": metrics.TotalOutputTokens,
			"rpm":                 metrics.RPM,
			"tpm":                 metrics.TPM,
			"average_response_ms": metrics.AverageResponseMs,
			"average_time_ms":     metrics.AverageTimeMs,
		}
		if p.paused[id] {
			pausedAccounts = append(pausedAccounts, id)
			continue
		}
		if p.inUse[id] < p.maxInflightPerAccount {
			available = append(available, id)
		}
	}
	for id, count := range p.inUse {
		if count > 0 {
			inUseAccounts = append(inUseAccounts, id)
			inUseSlots += count
		}
	}
	sort.Strings(inUseAccounts)
	sort.Strings(pausedAccounts)
	return map[string]any{
		"available":                len(available),
		"in_use":                   inUseSlots,
		"total":                    len(p.store.Accounts()),
		"paused":                   len(pausedAccounts),
		"available_accounts":       available,
		"in_use_accounts":          inUseAccounts,
		"paused_accounts":          pausedAccounts,
		"account_stats":            accountStats,
		"max_inflight_per_account": p.maxInflightPerAccount,
		"global_max_inflight":      p.globalMaxInflight,
		"recommended_concurrency":  p.recommendedConcurrency,
		"waiting":                  len(p.waiters),
		"max_queue_size":           p.maxQueueSize,
	}
}
