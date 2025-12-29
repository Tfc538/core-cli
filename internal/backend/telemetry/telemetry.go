package telemetry

import "context"

// Reporter is a placeholder for future telemetry integration.
type Reporter interface {
	Record(ctx context.Context, event string, fields map[string]any)
}
