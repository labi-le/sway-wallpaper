package wptool

import (
	"context"
	"os/exec"
)

type S3W struct{}

func (S3W) Set(ctx context.Context, path, output string) error {
	return exec.CommandContext(
		ctx,
		"swww",
		"img",
		path,
		"-o",
		output,
	).Start()
}

func (S3W) Clean() error {
	return exec.Command("swww", "clear").Start()
}
