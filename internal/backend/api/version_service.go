package api

import (
	"context"

	"github.com/Tfc538/core-cli/internal/backend/service/version"
)

// VersionService exposes version metadata for API handlers.
type VersionService interface {
	Latest(ctx context.Context) (version.Info, error)
	Get(ctx context.Context, version string) (version.Info, bool, error)
}
