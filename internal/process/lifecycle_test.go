package process

import (
	"context"
	"io"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/sthbryan/ftm/internal/config"
	"github.com/sthbryan/ftm/internal/providers"
)

// scriptProvider runs an arbitrary shell script as the "tunnel", so tests can
// make a provider crash, linger, or print a URL on demand.
type scriptProvider struct {
	script string
}

func (scriptProvider) Name() string       { return "Script" }
func (scriptProvider) BinaryName() string { return "sh" }

func (p scriptProvider) Start(ctx context.Context, _ config.TunnelConfig, logWriter io.Writer) (*providers.Process, error) {
	ctx, cancel := context.WithCancel(ctx)
	cmd := exec.CommandContext(ctx, "sh", "-c", p.script)
	cmd.Stdout = logWriter
	cmd.Stderr = logWriter

	return providers.StartProcess(cmd, cancel)
}

func (scriptProvider) ParseURL(line string) string {
	if len(line) > 8 && line[:8] == "https://" {
		return line
	}
	return ""
}

func newScriptManager(t *testing.T, script string) (*Manager, chan config.TunnelStatus) {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("uses a POSIX shell")
	}

	m := NewManager()
	m.providers[config.ProviderCloudflared] = scriptProvider{script: script}

	updates := make(chan config.TunnelStatus, 64)
	m.SetStatusChannel(updates)

	return m, updates
}

func tunnelFixture() config.TunnelConfig {
	return config.TunnelConfig{
		ID:        "t1",
		Name:      "Foundry VTT",
		Provider:  config.ProviderCloudflared,
		LocalPort: 30000,
	}
}

// waitForState drains updates until the wanted state shows up.
func waitForState(t *testing.T, updates <-chan config.TunnelStatus, want config.TunnelState) config.TunnelStatus {
	t.Helper()

	deadline := time.After(15 * time.Second)
	for {
		select {
		case status := <-updates:
			if status.State == want {
				return status
			}
		case <-deadline:
			t.Fatalf("never saw state %q", want)
		}
	}
}

// The bug: nothing waited on the process, so a provider that died kept being
// reported as online forever.
func TestManagerReportsCrashedProvider(t *testing.T) {
	m, updates := newScriptManager(t, "echo https://example.trycloudflare.com; sleep 0.2; exit 7")

	if err := m.Start(tunnelFixture(), nil); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	status := waitForState(t, updates, config.TunnelStateError)
	if status.ErrorMessage == "" {
		t.Error("error state carries no message")
	}
	if status.PublicURL != "" {
		t.Errorf("PublicURL = %q after the process died, want it cleared", status.PublicURL)
	}

	if _, ok := m.GetStatus("t1"); ok {
		t.Error("crashed tunnel is still registered as running")
	}
}

func TestManagerReportsCleanExitAsStopped(t *testing.T) {
	m, updates := newScriptManager(t, "exit 0")

	if err := m.Start(tunnelFixture(), nil); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	waitForState(t, updates, config.TunnelStateStopped)
}

// A deliberate Stop must not be reported as a crash.
func TestManagerStopIsNotReportedAsError(t *testing.T) {
	m, updates := newScriptManager(t, "sleep 60")

	if err := m.Start(tunnelFixture(), nil); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if err := m.Stop("t1"); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	waitForState(t, updates, config.TunnelStateStopped)

	// Nothing may follow claiming the tunnel failed.
	deadline := time.After(time.Second)
	for {
		select {
		case status := <-updates:
			if status.State == config.TunnelStateError {
				t.Fatalf("Stop produced an error state: %q", status.ErrorMessage)
			}
		case <-deadline:
			return
		}
	}
}

// Stop has to return only once the process is really gone, otherwise quitting
// leaves tunnels running.
func TestManagerStopWaitsForExit(t *testing.T) {
	m, _ := newScriptManager(t, "sleep 60")

	if err := m.Start(tunnelFixture(), nil); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	m.mu.RLock()
	proc := m.processes["t1"].Process
	m.mu.RUnlock()

	if err := m.Stop("t1"); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	select {
	case <-proc.Exited():
	default:
		t.Fatal("Stop returned while the process was still running")
	}
}

func TestManagerStopAllWaitsForEveryProcess(t *testing.T) {
	m, _ := newScriptManager(t, "sleep 60")

	for _, id := range []string{"t1", "t2", "t3"} {
		tunnel := tunnelFixture()
		tunnel.ID = id
		if err := m.Start(tunnel, nil); err != nil {
			t.Fatalf("Start(%s) failed: %v", id, err)
		}
	}

	m.mu.RLock()
	procs := make([]*providers.Process, 0, len(m.processes))
	for _, mp := range m.processes {
		procs = append(procs, mp.Process)
	}
	m.mu.RUnlock()

	if len(procs) != 3 {
		t.Fatalf("started %d processes, want 3", len(procs))
	}

	m.StopAll()

	for i, proc := range procs {
		select {
		case <-proc.Exited():
		default:
			t.Fatalf("process %d still running after StopAll", i)
		}
	}
}

// Stopping and restarting inside the startup window used to let the old
// monitor time out the new process.
func TestManagerRestartIsNotTimedOutByTheOldMonitor(t *testing.T) {
	m, updates := newScriptManager(t, "sleep 60")

	if err := m.Start(tunnelFixture(), nil); err != nil {
		t.Fatalf("first Start failed: %v", err)
	}
	if err := m.Stop("t1"); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
	if err := m.Start(tunnelFixture(), nil); err != nil {
		t.Fatalf("second Start failed: %v", err)
	}

	// The first monitor's "connecting" step lands ~5s in; give it room to fire.
	deadline := time.After(9 * time.Second)
	for {
		select {
		case status := <-updates:
			if status.State == config.TunnelStateTimeout {
				t.Fatal("the stale monitor timed out the restarted tunnel")
			}
		case <-deadline:
			if _, ok := m.GetStatus("t1"); !ok {
				t.Fatal("restarted tunnel is no longer registered")
			}
			m.StopAll()
			return
		}
	}
}
