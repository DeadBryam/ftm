//go:build darwin

package notifications

func newPlatformNotifier() (Notifier, SoundPlayer) {
	return newDarwinNotifier(), newDarwinSoundPlayer()
}
