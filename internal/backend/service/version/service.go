package version

import "context"

// Provider returns version metadata from a backing store.
type Provider interface {
	Latest(ctx context.Context) (Info, error)
	Get(ctx context.Context, version string) (Info, bool, error)
}

// Service exposes version metadata for API handlers.
type Service struct {
	provider Provider
}

// NewService constructs a version service with the given provider.
func NewService(provider Provider) *Service {
	return &Service{provider: provider}
}

// Latest returns the most recent version metadata.
func (s *Service) Latest(ctx context.Context) (Info, error) {
	return s.provider.Latest(ctx)
}

// Get returns metadata for a specific version, if available.
func (s *Service) Get(ctx context.Context, version string) (Info, bool, error) {
	return s.provider.Get(ctx, version)
}
