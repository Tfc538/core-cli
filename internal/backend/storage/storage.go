package storage

import (
	"context"
	"io"
)

// BinaryStore abstracts access to versioned binary artifacts.
type BinaryStore interface {
	OpenBinary(ctx context.Context, version string, platform string) (io.ReadCloser, error)
}

// TemplateStore abstracts access to versioned templates.
type TemplateStore interface {
	GetTemplate(ctx context.Context, name string, version string) ([]byte, error)
}
