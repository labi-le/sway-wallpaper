package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/labi-le/sway-wallpaper/internal/manager"
	"github.com/nightlyone/lockfile"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
	"os"
	"os/signal"
	"path/filepath"
)

var (
	ErrSwayWallpaperAlreadyRunning = errors.New("sway-wallpaper is already running. pid %d")
	ErrCannotLock                  = errors.New("cannot get locked process: %s")
	ErrCannotUnlock                = errors.New("cannot unlock process: %s")
)

const LockFile = "sway-wallpaper.lck"

func main() {
	lock := MustLock()
	defer Unlock(lock)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opt, err := manager.ParseOptions()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	wp, err := manager.New(opt)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer wp.Close()

	initLogger(opt.Debug)

	if opt.FollowDuration == 0 {
		if wpErr := wp.Provide(ctx); wpErr != nil {
			if errors.Is(wpErr, context.Canceled) {
				log.Warn().Msg("Handling interrupt signal")
			}

			return
		}
		<-ctx.Done()
		return
	}

	for {
		reuse, currentCancel := context.WithTimeout(ctx, opt.FollowDuration)

		if wpErr := wp.Provide(reuse); wpErr != nil {
			if errors.Is(wpErr, context.Canceled) {
				log.Warn().Msg("Handling interrupt signal")
				break
			}
			log.Fatal().Err(wpErr).Msg("start error")
		}

		<-reuse.Done()
		currentCancel()
	}
}

func MustLock() lockfile.Lockfile {
	lock, _ := lockfile.New(filepath.Join(os.TempDir(), LockFile))

	if lockErr := lock.TryLock(); lockErr != nil {
		owner, err := lock.GetOwner()
		if err != nil {
			log.Fatal().Err(fmt.Errorf("%w: %v", ErrCannotLock, err))
		}
		log.Fatal().Err(fmt.Errorf("%w: pid: %d", ErrSwayWallpaperAlreadyRunning, owner.Pid))
	}

	return lock
}

func Unlock(lock lockfile.Lockfile) {
	if err := lock.Unlock(); err != nil {
		log.Fatal().Err(fmt.Errorf("%w: %v", ErrCannotUnlock, err))
	}
}

func initLogger(debug bool) {
	if debug {
		log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})

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
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		return
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
