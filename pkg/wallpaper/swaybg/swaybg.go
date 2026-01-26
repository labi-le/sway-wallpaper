package swaybg

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/labi-le/chiasma/pkg/wallpaper/execute"
)

const Name = "swaybg"

type SwayBG struct{}

func NewSwayBG() (SwayBG, error) {
	if _, err := exec.LookPath(Name); err != nil {
		return SwayBG{}, fmt.Errorf("%s: %w", Name, execute.ErrUtilityNotFound)
	}

	return SwayBG{}, nil
}

func (t SwayBG) Change(ctx context.Context, path, output string) error {
	return exec.CommandContext(
		ctx,
		Name,
		"-i",
		path,
		"-o",
		output,
	).Start()
}
