package tools

import (
	"golang.org/x/net/context"
	"os/exec"
)

func SetWallpaperWBG(ctx context.Context, path string) error {
	cmd := exec.CommandContext(ctx, "wbg", path)
	return cmd.Start()
}
