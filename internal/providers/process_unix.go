//go:build !windows

package providers

import (
	"os/exec"
	"syscall"
)

// configureProcessGroup puts the child in its own process group so the whole
// tree can be signalled at once.
//
// ssh and cloudflared spawn helpers of their own. Cancelling the context only
// kills the direct child, so without this those helpers survive and the tunnel
// stays up after ftm exits -- leaving the user's Foundry world exposed by a
// tunnel they believe they closed.
func configureProcessGroup(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
}

// signalGroup sends sig to the child's entire process group, falling back to
// the child alone if the group cannot be resolved.
func signalGroup(cmd *exec.Cmd, sig syscall.Signal) {
	if cmd == nil || cmd.Process == nil {
		return
	}

	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		_ = cmd.Process.Signal(sig)
		return
	}

	// A negative pid addresses the whole group.
	_ = syscall.Kill(-pgid, sig)
}

func terminateGroup(cmd *exec.Cmd) {
	signalGroup(cmd, syscall.SIGTERM)
}

func killGroup(cmd *exec.Cmd) {
	signalGroup(cmd, syscall.SIGKILL)
}
