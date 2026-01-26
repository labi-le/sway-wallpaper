package swww

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/labi-le/chiasma/pkg/wallpaper/execute"
)

const Name = "swww"

type SWWW struct{}

func NewSWWW() (SWWW, error) {
	if _, err := exec.LookPath(Name); err != nil {
		return SWWW{}, fmt.Errorf("%s: %w", Name, execute.ErrUtilityNotFound)
	}

	return SWWW{}, nil
}

func (t SWWW) Change(ctx context.Context, path, output string) error {
	return exec.CommandContext(
		ctx,
		Name,
		"img",
		path,
		"-o",
		output,
	).Start()
}

func (t SWWW) Close() error {
	return exec.Command(Name, "clear").Start()
}
