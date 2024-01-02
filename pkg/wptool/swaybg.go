package wptool

import (
	"context"
	"os/exec"
)

type SwayBG struct{}

func (SwayBG) Set(ctx context.Context, path, output string) error {
	return exec.CommandContext(
		ctx,
		"swaybg",
		"-i",
		path,
		"-o",
		output,
	).Start()
}
