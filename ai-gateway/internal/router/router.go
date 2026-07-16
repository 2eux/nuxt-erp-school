package router

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/opencode/erp-ai-gateway/internal/config"
	"github.com/opencode/erp-ai-gateway/internal/providers"
)

type ProviderRouter struct {
	cfg       *config.Config
	providers map[string]providers.Provider
	mu        sync.RWMutex
	counters  map[string]*int64
}

func NewProviderRouter(cfg *config.Config, providerMap map[string]providers.Provider) *ProviderRouter {
	counters := make(map[string]*int64)
	for name := range providerMap {
		var c int64
		counters[name] = &c
	}

	return &ProviderRouter{
		cfg:       cfg,
		providers: providerMap,
		counters:  counters,
	}
}

func (r *ProviderRouter) GetProvider(name string) (providers.Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not available", name)
	}
	return p, nil
}

func (r *ProviderRouter) GetAllProviders() []providers.Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]providers.Provider, 0)
	for _, p := range r.providers {
		result = append(result, p)
	}
	return result
}

func (r *ProviderRouter) GetEnabledProviders() []providers.Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]providers.Provider, 0)
	for _, name := range r.cfg.Enabled {
		if p, ok := r.providers[name]; ok {
			result = append(result, p)
		}
	}
	return result
}

func (r *ProviderRouter) Route(ctx context.Context, preferredProvider string, taskType string) (providers.Provider, error) {
	if preferredProvider != "" {
		p, err := r.GetProvider(preferredProvider)
		if err == nil && r.isHealthy(ctx, p) {
			r.incrementCounter(p.Name())
			return p, nil
		}
	}

	switch taskType {
	case "embedding":
		return r.routeForEmbedding(ctx)
	case "streaming":
		return r.routeLoadBalanced(ctx)
	case "islamic":
		return r.routeForIslamic(ctx)
	case "cheap":
		return r.routeCheapest(ctx)
	default:
		return r.routeWithFallback(ctx)
	}
}

func (r *ProviderRouter) routeForEmbedding(ctx context.Context) (providers.Provider, error) {
	preferredOrder := []string{"openai", "gemini", "ollama"}

	for _, name := range preferredOrder {
		if !r.cfg.IsProviderEnabled(name) {
			continue
		}
		p, err := r.GetProvider(name)
		if err != nil {
			continue
		}
		if r.isHealthy(ctx, p) {
			r.incrementCounter(p.Name())
			return p, nil
		}
	}

	return nil, fmt.Errorf("no embedding provider available")
}

func (r *ProviderRouter) routeLoadBalanced(ctx context.Context) (providers.Provider, error) {
	enabled := r.GetEnabledProviders()
	healthy := make([]providers.Provider, 0)
	for _, p := range enabled {
		if r.isHealthy(ctx, p) {
			healthy = append(healthy, p)
		}
	}

	if len(healthy) == 0 {
		return nil, fmt.Errorf("no healthy providers available")
	}

	idx := rand.Intn(len(healthy))
	r.incrementCounter(healthy[idx].Name())
	return healthy[idx], nil
}

func (r *ProviderRouter) routeWithFallback(ctx context.Context) (providers.Provider, error) {
	for _, name := range r.cfg.FallbackChain {
		if !r.cfg.IsProviderEnabled(name) {
			continue
		}
		p, err := r.GetProvider(name)
		if err != nil {
			continue
		}
		if r.isHealthy(ctx, p) {
			r.incrementCounter(p.Name())
			return p, nil
		}
	}

	defaultProvider, err := r.GetProvider(r.cfg.DefaultProvider)
	if err == nil && r.isHealthy(ctx, defaultProvider) {
		r.incrementCounter(defaultProvider.Name())
		return defaultProvider, nil
	}

	enabled := r.GetEnabledProviders()
	for _, p := range enabled {
		if r.isHealthy(ctx, p) {
			r.incrementCounter(p.Name())
			return p, nil
		}
	}

	return nil, fmt.Errorf("all providers are unavailable")
}

func (r *ProviderRouter) routeForIslamic(ctx context.Context) (providers.Provider, error) {
	if p, err := r.GetProvider("ollama"); err == nil && r.isHealthy(ctx, p) {
		r.incrementCounter(p.Name())
		return p, nil
	}

	for _, name := range r.cfg.FallbackChain {
		if !r.cfg.IsProviderEnabled(name) {
			continue
		}
		p, err := r.GetProvider(name)
		if err != nil {
			continue
		}
		if r.isHealthy(ctx, p) {
			r.incrementCounter(p.Name())
			return p, nil
		}
	}

	return nil, fmt.Errorf("no provider available for islamic tasks")
}

func (r *ProviderRouter) routeCheapest(ctx context.Context) (providers.Provider, error) {
	cheapest := struct {
		name string
		cost float64
	}{}

	for name, cfg := range r.cfg.Providers {
		if !r.cfg.IsProviderEnabled(name) {
			continue
		}
		totalCost := cfg.CostPer1KIn + cfg.CostPer1KOut
		if cheapest.name == "" || totalCost < cheapest.cost {
			if p, err := r.GetProvider(name); err == nil && r.isHealthy(ctx, p) {
				cheapest.name = name
				cheapest.cost = totalCost
			}
		}
	}

	if cheapest.name == "" {
		return r.routeWithFallback(ctx)
	}

	p, _ := r.GetProvider(cheapest.name)
	r.incrementCounter(p.Name())
	return p, nil
}

func (r *ProviderRouter) isHealthy(ctx context.Context, p providers.Provider) bool {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return p.IsAvailable(ctx)
}

func (r *ProviderRouter) incrementCounter(name string) {
	if c, ok := r.counters[name]; ok {
		atomic.AddInt64(c, 1)
	}
}

func (r *ProviderRouter) GetStats() map[string]int64 {
	stats := make(map[string]int64)
	r.mu.RLock()
	defer r.mu.RUnlock()
	for name, c := range r.counters {
		stats[name] = atomic.LoadInt64(c)
	}
	return stats
}
