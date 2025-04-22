package wallpaper

import (
	"context"
	"fmt"
	"os/exec"
)

type Tool interface {
	Change(ctx context.Context, path, output string) error
	String() string
}

type ToolDestructor interface {
	Close() error
}

// Stringer is an interface for types that can convert themselves to a string
type Stringer interface {
	String() string
}

// lookup is a generic function that checks if a binary is in PATH
func lookup[T Stringer](t T) (T, error) {
	if _, err := exec.LookPath(t.String()); err != nil {
		return t, fmt.Errorf("%s not found in PATH", t)
	}

	return t, nil
}
