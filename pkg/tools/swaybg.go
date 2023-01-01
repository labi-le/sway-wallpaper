package tools

import "os/exec"

func SetWallpaperSwayBG(path string) error {
	_ = exec.Command("killall", "swaybg").Run()
	// since this process does not end, we ignore it
	_ = exec.Command("swaybg", "-i", path).Start()

	return nil
}
