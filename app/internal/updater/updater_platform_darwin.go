//go:build darwin

package updater

import (
	"os"
	"syscall"
)

func init() {
	platformUpdaterFactory = func() PlatformUpdater {
		return &darwinPlatformUpdater{}
	}
}

type darwinPlatformUpdater struct{}

func (d *darwinPlatformUpdater) CanReplaceRunningExecutable() bool {
	// macOS allows replacing a running executable
	return true
}

func (d *darwinPlatformUpdater) ReplaceExecutable(src, dst string) error {
	// On macOS, we can directly replace the file
	return os.Rename(src, dst)
}

func (d *darwinPlatformUpdater) MakeExecutable(path string) error {
	return os.Chmod(path, 0o755)
}

func (d *darwinPlatformUpdater) RestartApplication() error {
	// Get current executable path
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	// Prepare command arguments (same as current process)
	args := os.Args

	// Use exec to replace the current process
	return syscall.Exec(exe, args, os.Environ())
}
