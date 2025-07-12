//go:build darwin

package updater

import "os"

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
