//go:build linux

package updater

import (
	"os"
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
