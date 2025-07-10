//go:build windows

package updater

import (
	"os"
	"time"
)

func init() {
	platformUpdaterFactory = func() PlatformUpdater {
		return &windowsPlatformUpdater{}
	}
}

type windowsPlatformUpdater struct{}

func (w *windowsPlatformUpdater) CanReplaceRunningExecutable() bool {
	// Windows doesn't allow replacing a running executable
	return false
}

func (w *windowsPlatformUpdater) ReplaceExecutable(src, dst string) error {
	// Rename current executable
	oldPath := dst + ".old"
	if err := os.Rename(dst, oldPath); err != nil {
		return err
	}
	
	// Move new executable
	if err := os.Rename(src, dst); err != nil {
		// Restore old executable
		os.Rename(oldPath, dst)
		return err
	}
	
	// Schedule old executable for deletion
	// Windows will delete it after the process exits
	go func() {
		time.Sleep(5 * time.Second)
		os.Remove(oldPath)
	}()
	
	return nil
}

func (w *windowsPlatformUpdater) MakeExecutable(path string) error {
	// Windows doesn't need chmod
	return nil
}