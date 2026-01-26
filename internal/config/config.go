package config

import (
	"os"
	"time"

	output2 "github.com/labi-le/chiasma/internal/output"
	"github.com/labi-le/chiasma/pkg/api"
	"github.com/labi-le/chiasma/pkg/browser"
	flag "github.com/spf13/pflag"
)

type Config struct {
	BrowserName    string
	HistoryPath    string
	Resolution     output2.Resolution
	OutputMonitor  output2.Monitor
	ToolName       string
	APIName        string
	SaveDir        string
	SearchPhrase   string
	Follow         bool
	FollowDuration time.Duration
	Verbose        bool
}

func Parse() (Config, error) {
	var c Config
	flag.StringVar(&c.BrowserName, "browser", browser.AvailableBrowsers()[0], "browser name")
	flag.StringVar(&c.HistoryPath, "history-file", "", "path to history file")
	flag.Var(&c.Resolution, "resolution", "target resolution (e.g. 1920x1080)")
	flag.Var(&c.OutputMonitor, "output", "monitor output (e.g. eDP-1)")
	flag.StringVar(&c.ToolName, "tool", "", "wallpaper tool")
	flag.StringVar(&c.APIName, "api", api.AvailableAPIs()[0], "image source api")
	flag.StringVar(&c.SaveDir, "save-dir", os.Getenv("HOME")+"/Pictures/chiasma", "save directory")
	flag.StringVar(&c.SearchPhrase, "phrase", "", "search phrase")
	flag.DurationVar(&c.FollowDuration, "interval", time.Hour, "update interval")
	flag.BoolVar(&c.Follow, "follow", false, "enable periodic updates")
	flag.BoolVar(&c.Verbose, "verbose", false, "enable verbose logs")

	flag.Parse()

	return c, nil
}
