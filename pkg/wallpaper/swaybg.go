package wallpaper

import (
	"context"
	"os/exec"
)

type SwayBG struct{}

func NewSwayBG() (SwayBG, error) {
	return lookup(SwayBG{})
}

func (t SwayBG) Change(ctx context.Context, path, output string) error {
	return exec.CommandContext(
		ctx,
		t.String(),
		"-i",
		path,
		"-o",
		output,
	).Start()
}

func (t SwayBG) String() string {
	return "swaybg"
}
