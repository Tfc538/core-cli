package version

import (
	"context"
	"errors"
	"sync"
)

var errNoVersions = errors.New("no versions available")

// InMemoryProvider stores version metadata in memory.
type InMemoryProvider struct {
	mu      sync.RWMutex
	latest  string
	entries map[string]Info
}

// NewInMemoryProvider builds an in-memory provider from the supplied list.
func NewInMemoryProvider(entries []Info) *InMemoryProvider {
	provider := &InMemoryProvider{
		entries: make(map[string]Info),
	}

	for _, entry := range entries {
		provider.entries[entry.Version] = entry
		provider.latest = entry.Version
	}

	return provider
}

// Latest returns the most recent version metadata.
func (p *InMemoryProvider) Latest(ctx context.Context) (Info, error) {
	_ = ctx

	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.latest == "" {
		return Info{}, errNoVersions
	}

	entry, ok := p.entries[p.latest]
	if !ok {
		return Info{}, errNoVersions
	}

	return entry, nil
}

// Get returns metadata for a specific version.
func (p *InMemoryProvider) Get(ctx context.Context, version string) (Info, bool, error) {
	_ = ctx

	p.mu.RLock()
	defer p.mu.RUnlock()

	entry, ok := p.entries[version]
	if !ok {
		return Info{}, false, nil
	}

	return entry, true, nil
}
