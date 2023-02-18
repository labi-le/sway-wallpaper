package wptool

import (
	"golang.org/x/net/context"
	"os/exec"
)

type WBG struct{}

func (WBG) Set(ctx context.Context, path string) error {
	cmd := exec.CommandContext(ctx, "wbg", path)
	return cmd.Start()
}
