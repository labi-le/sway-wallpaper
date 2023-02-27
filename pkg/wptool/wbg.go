package wptool

import (
	"golang.org/x/net/context"
	"os/exec"
)

type WBG struct{}

func (WBG) Set(ctx context.Context, path string) error {
	return exec.CommandContext(ctx, "wbg", path).Start()
}
