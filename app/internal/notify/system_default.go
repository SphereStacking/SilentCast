//go:build !windows && !darwin && !linux

package notify

func init() {
	// Default implementation for unsupported platforms
	getSystemNotifier = func() Notifier {
		return nil
	}
}
