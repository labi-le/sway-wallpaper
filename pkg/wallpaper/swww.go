package wallpaper

import (
	"context"
	"os/exec"
)

type SWWW struct{}

func NewSWWW() (SWWW, error) {
	return lookup(SWWW{})
}

func (t SWWW) Change(ctx context.Context, path, output string) error {
	return exec.CommandContext(
		ctx,
		t.String(),
		"img",
		path,
		"-o",
		output,
	).Start()
}

func (t SWWW) Close() error {
	return exec.Command(t.String(), "clear").Start()
}

func (t SWWW) String() string {
	return "swww"
}
