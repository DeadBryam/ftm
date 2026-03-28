//go:build linux

package notifications

func newPlatformNotifier() (Notifier, SoundPlayer) {
	return newLinuxNotifier(), newLinuxSoundPlayer()
}
