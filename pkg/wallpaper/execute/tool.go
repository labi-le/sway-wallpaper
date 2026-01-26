package execute

import (
	"context"
	"errors"
)

var (
	ErrUtilityNotFound = errors.New("utility not found in PATH")
)

type Provider interface {
	Change(ctx context.Context, path, output string) error
}
