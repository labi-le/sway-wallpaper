package main

import (
	"context"
	"errors"
	"github.com/labi-le/sway-wallpaper/internal/manager"
	"github.com/labi-le/sway-wallpaper/pkg/log"
	"github.com/nightlyone/lockfile"
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

	opt := manager.ParseOptions()
	wp := manager.New(opt)

	if opt.FollowDuration == 0 {
		if wpErr := wp.Set(ctx); wpErr != nil {
			if errors.Is(wpErr, context.Canceled) {
				log.Warnf("Handling interrupt signal")
				return
			}
			log.Error(wpErr)
		}
		<-ctx.Done()
		return
	}

	for {
		reuse, currentCancel := context.WithTimeout(ctx, opt.FollowDuration)

		if wpErr := wp.Set(reuse); wpErr != nil {
			if errors.Is(wpErr, context.Canceled) {
				log.Warnf("Handling interrupt signal")
				break
			}
			log.Error(wpErr)
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
			log.Fatalf(ErrCannotLock.Error(), err)
		}
		log.Fatalf(ErrSwayWallpaperAlreadyRunning.Error(), owner.Pid)
	}

	return lock
}

func Unlock(lock lockfile.Lockfile) {
	if err := lock.Unlock(); err != nil {
		log.Fatalf(ErrCannotUnlock.Error(), err)
	}
}
