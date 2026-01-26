package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labi-le/chiasma/internal/config"
	"github.com/labi-le/chiasma/internal/output"
	"github.com/labi-le/chiasma/internal/service"
	"github.com/labi-le/chiasma/pkg/api"
	"github.com/labi-le/chiasma/pkg/browser"
	"github.com/labi-le/chiasma/pkg/wallpaper"
	"github.com/rs/zerolog"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		panic(err)
	}

	log := initLogger(cfg.Verbose)

	searcher, err := api.NewSearcher(log, cfg.APIName)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init api")
	}

	var historyProvider service.QuerySource
	if cfg.SearchPhrase == "" {
		hp, err := browser.NewHistoryProvider(cfg.BrowserName, cfg.HistoryPath)
		if err != nil {
			log.Warn().Err(err).Msg("failed to init browser history, fallback to random or manual phrase might fail")
		} else {
			defer func() {
				if closer, ok := hp.(interface{ Close() error }); ok {
					closer.Close()
				}
			}()
			historyProvider = hp
		}
	}

	tool, err := wallpaper.ByNameOrAvailable(cfg.ToolName)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init wallpaper tool")
	}

	resolution := cfg.Resolution
	if resolution.Width == 0 || resolution.Height == 0 {
		mon, err := output.NewByIDXrandr(cfg.OutputMonitor.ID)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to detect resolution, please specify --image-resolution")
		}
		resolution = mon.CurrentResolution
		log.Info().Str("res", resolution.String()).Msg("detected resolution")
	}

	svc := &service.WallpaperService{
		Log:     log,
		API:     searcher,
		History: historyProvider,
		Setter:  tool,
	}

	params := service.UpdateParams{
		Phrase:     cfg.SearchPhrase,
		Resolution: resolution,
		SaveDir:    cfg.SaveDir,
		OutputID:   cfg.OutputMonitor.ID,
		RetryCount: 5,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	run := func() {
		if err := svc.Update(ctx, params); err != nil {
			log.Error().Err(err).Msg("failed to update wallpaper")
		}
	}

	run()

	if !cfg.Follow {
		return
	}

	ticker := time.NewTicker(cfg.FollowDuration)
	defer ticker.Stop()

	log.Info().Dur("interval", cfg.FollowDuration).Msg("entering watch mode")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("shutting down")
			return
		case <-ticker.C:
			run()
		}
	}
}

func initLogger(verbose bool) zerolog.Logger {
	out := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}

	if verbose {
		zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
			return fmt.Sprintf("%s:%d", file, line)
		}
		return zerolog.New(out).
			Level(zerolog.TraceLevel).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	return zerolog.New(out).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
}
