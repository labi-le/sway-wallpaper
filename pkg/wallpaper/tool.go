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

type Stringer interface {
	String() string
}

func lookup[T Stringer](t T) (T, error) {
	if _, err := exec.LookPath(t.String()); err != nil {
		return t, fmt.Errorf("%s not found in PATH", t)
	}

	return t, nil
}
