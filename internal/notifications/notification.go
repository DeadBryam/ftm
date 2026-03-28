package notifications

import (
	"fmt"
	"log"
)

type Notifier interface {
	Notify(title, body string) error
	IsAvailable() bool
}

type SoundPlayer interface {
	PlaySound(t SoundType) error
	IsAvailable() bool
}

type SoundType int

const (
	SoundStartup SoundType = iota
	SoundSuccess
	SoundError
	SoundWarning
	SoundAlert
	SoundInfo
)

func (s SoundType) String() string {
	switch s {
	case SoundStartup:
		return "startup"
	case SoundSuccess:
		return "success"
	case SoundError:
		return "error"
	case SoundWarning:
		return "warning"
	case SoundAlert:
		return "alert"
	case SoundInfo:
		return "info"
	default:
		return "unknown"
	}
}

var (
	notifier    Notifier
	soundPlayer SoundPlayer
	available   bool
)

func Init() {
	n, s := newPlatformNotifier()
	if n != nil && n.IsAvailable() {
		notifier = n
		soundPlayer = s
		available = true
		log.Printf("[notifications] Using platform-native")
		return
	}

	n, s = newGoNotification()
	if n != nil && n.IsAvailable() {
		notifier = n
		soundPlayer = s
		available = true
		log.Printf("[notifications] Using fallback")
		return
	}

	available = false
	log.Printf("[notifications] Running in silent mode")
}

func Notify(title, body string) error {
	if notifier == nil || !available {
		log.Printf("[notification] %s: %s", title, body)
		return nil
	}
	return notifier.Notify(title, body)
}

func Notifyf(title, format string, args ...interface{}) error {
	return Notify(title, fmt.Sprintf(format, args...))
}

func PlaySound(t SoundType) error {
	if soundPlayer == nil || !available {
		return nil
	}
	return soundPlayer.PlaySound(t)
}

func IsAvailable() bool {
	return available
}

func NotifyTunnelOnline(name, url string) {
	Notify("Tunnel Active", fmt.Sprintf("%s - %s", name, url))
	PlaySound(SoundSuccess)
}

func NotifyTunnelError(name, errMsg string) {
	Notify("Tunnel Error", fmt.Sprintf("%s: %s", name, errMsg))
	PlaySound(SoundError)
}

func NotifyTunnelTimeout(name string) {
	Notify("Tunnel Timeout", fmt.Sprintf("%s could not connect", name))
	PlaySound(SoundError)
}

func NotifyTunnelStopped(name string) {
	Notify("Tunnel Stopped", fmt.Sprintf("%s has been stopped", name))
	PlaySound(SoundInfo)
}

func NotifyTunnelExpiring(name string, minutes int) {
	if minutes == 1 {
		Notify("Last Minute!", fmt.Sprintf("%s: 1 minute remaining", name))
		PlaySound(SoundAlert)
	} else {
		Notify("Tunnel Expiring", fmt.Sprintf("%s: %d minutes remaining", name, minutes))
		PlaySound(SoundWarning)
	}
}

func NotifyTunnelExpired(name string) {
	Notify("Tunnel Expired", fmt.Sprintf("%s session has ended", name))
	PlaySound(SoundError)
}

func NotifyWelcome() {
	Notify("Welcome!", "Foundry Tunnel Manager started")
	PlaySound(SoundStartup)
}
