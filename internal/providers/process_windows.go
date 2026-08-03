//go:build windows

package providers

import (
	"os/exec"
	"strconv"
	"syscall"
)

// configureProcessGroup puts the child in a new process group, so that a Ctrl+C
// delivered to ftm's own console is not broadcast to the tunnel as well.
func configureProcessGroup(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.CreationFlags |= syscall.CREATE_NEW_PROCESS_GROUP
}

// taskkill ends the process tree rooted at the child.
//
// Windows has no signals, and Process.Kill reaches only the direct child, so an
// ssh.exe or cloudflared.exe helper would otherwise survive and keep the tunnel
// open. /T includes descendants; /F makes it non-negotiable.
func taskkill(cmd *exec.Cmd, force bool) {
	if cmd == nil || cmd.Process == nil {
		return
	}

	args := []string{"/T", "/PID", strconv.Itoa(cmd.Process.Pid)}
	if force {
		args = append([]string{"/F"}, args...)
	}

	kill := exec.Command("taskkill", args...)
	kill.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_ = kill.Run()
}

func terminateGroup(cmd *exec.Cmd) {
	taskkill(cmd, false)
}

func killGroup(cmd *exec.Cmd) {
	taskkill(cmd, true)
}
