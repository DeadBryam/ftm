//go:build !darwin

package notifications

import "runtime"

func newPlatformNotifier() (Notifier, SoundPlayer) {
	switch runtime.GOOS {
	case "linux":
		return newLinuxNotifier(), newLinuxSoundPlayer()
	case "windows":
		return newWindowsNotifier(), newWindowsSoundPlayer()
	}
	return nil, nil
}
