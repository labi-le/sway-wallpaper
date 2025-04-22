package manager

import (
	"context"
)

type Setter interface {
	Change(ctx context.Context, path, output string) error
	Close() error
}
