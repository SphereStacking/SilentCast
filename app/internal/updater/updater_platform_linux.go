//go:build linux

package updater

import (
	"os"
	"syscall"
)

func init() {
	platformUpdaterFactory = func() PlatformUpdater {
		return &linuxUpdater{}
	}
}

type linuxUpdater struct{}

func (l *linuxUpdater) CanReplaceRunningExecutable() bool {
	// Linux allows replacing running executables
	return true
}

func (l *linuxUpdater) ReplaceExecutable(src, dst string) error {
	// On Linux, we can directly replace the executable
	return os.Rename(src, dst)
}

func (l *linuxUpdater) MakeExecutable(path string) error {
	// Set executable permissions (755)
	return os.Chmod(path, 0o755)
}

func (l *linuxUpdater) RestartApplication() error {
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
