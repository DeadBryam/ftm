package buildinfo

import "sync/atomic"

var desktop atomic.Bool

func MarkDesktop() {
	desktop.Store(true)
}

func IsDesktop() bool {
	return desktop.Load()
}
