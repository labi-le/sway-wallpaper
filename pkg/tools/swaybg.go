package tools

import (
	"golang.org/x/net/context"
	"os/exec"
)

func SetWallpaperSwayBG(ctx context.Context, path string) error {
	cmd := exec.CommandContext(ctx, "swaybg", "-i", path)
	return cmd.Start()
}
