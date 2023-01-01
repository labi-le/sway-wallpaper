package tools

import "os/exec"

func SetWallpaperWBG(path string) error {
	_ = exec.Command("killall", "wbg").Run()
	// since this process does not end, we ignore it
	_ = exec.Command("wbg", path).Start()

	return nil
}
