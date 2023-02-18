package wptool

import (
	"golang.org/x/net/context"
	"os/exec"
)

type SwayBG struct{}

func (SwayBG) Set(ctx context.Context, path string) error {
	cmd := exec.CommandContext(ctx, "swaybg", "-i", path)
	return cmd.Start()
}
