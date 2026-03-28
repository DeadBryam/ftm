//go:build windows

package notifications

func newPlatformNotifier() (Notifier, SoundPlayer) {
	return newWindowsNotifier(), newWindowsSoundPlayer()
}
