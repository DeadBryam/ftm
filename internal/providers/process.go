package providers

import (
	"context"
	"os/exec"
	"sync"
	"time"
)

// stopGrace is how long a process gets to shut down cleanly after being asked
// before it is killed outright.
const stopGrace = 3 * time.Second

// Process is a running tunnel provider.
//
// It owns the exec.Cmd rather than just a CancelFunc. Without the Cmd nothing
// ever calls Wait, which had two consequences: the child was never reaped, and
// a provider that died on its own went unnoticed, leaving the UI showing
// "online" for a tunnel that no longer existed.
type Process struct {
	Cmd       *exec.Cmd
	Cancel    context.CancelFunc
	PublicURL string

	exited   chan struct{}
	waitErr  error
	stopOnce sync.Once
}

// StartProcess starts cmd with the setup needed to tear down its whole process
// tree, and begins reaping it in the background.
//
// cancel must be the CancelFunc of the context cmd was built with; it is called
// if the start fails and again on Stop.
func StartProcess(cmd *exec.Cmd, cancel context.CancelFunc) (*Process, error) {
	configureProcessGroup(cmd)

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, err
	}

	p := &Process{
		Cmd:    cmd,
		Cancel: cancel,
		exited: make(chan struct{}),
	}

	go func() {
		p.waitErr = cmd.Wait()
		close(p.exited)
	}()

	return p, nil
}

// Exited is closed once the process has terminated and been reaped.
func (p *Process) Exited() <-chan struct{} {
	return p.exited
}

// Err reports why the process exited. Only meaningful once Exited is closed;
// nil means it exited successfully.
func (p *Process) Err() error {
	return p.waitErr
}

// Stop terminates the process and everything it spawned, and waits for it to
// go away.
//
// Callers must not hold a lock across this: it blocks for up to two grace
// periods in the worst case.
func (p *Process) Stop() {
	p.stopOnce.Do(func() {
		terminateGroup(p.Cmd)

		select {
		case <-p.exited:
		case <-time.After(stopGrace):
			// Ignored the polite request; take the tree down hard.
			killGroup(p.Cmd)
			select {
			case <-p.exited:
			case <-time.After(stopGrace):
			}
		}

		if p.Cancel != nil {
			p.Cancel()
		}
	})
}
