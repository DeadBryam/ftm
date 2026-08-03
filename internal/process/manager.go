package process

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/sthbryan/ftm/internal/config"
	"github.com/sthbryan/ftm/internal/providers"
	"github.com/sthbryan/ftm/internal/providers/bore"
	"github.com/sthbryan/ftm/internal/providers/cloudflared"
	"github.com/sthbryan/ftm/internal/providers/pinggy"
	"github.com/sthbryan/ftm/internal/providers/ssh"
	"github.com/sthbryan/ftm/internal/providers/tunnelmole"
)

type Manager struct {
	mu                  sync.RWMutex
	processes           map[string]*ManagedProcess
	providers           map[config.Provider]providers.Provider
	providerExpiration  map[string]int
	DownloadProgress    chan providers.DownloadProgress
	StatusChannel       chan config.TunnelStatus
	NotificationHandler func(status config.TunnelStatus)
	ExpirationCallbacks struct {
		OnStart func(tunnelID, name, provider string, startedAt time.Time)
		OnStop  func(tunnelID string)
	}
}

func (m *Manager) SetProgressChannel(ch chan providers.DownloadProgress) {
	m.DownloadProgress = ch
	for _, p := range m.providers {
		if installer, ok := p.(interface {
			SetProgressChannel(chan providers.DownloadProgress)
		}); ok {
			installer.SetProgressChannel(ch)
		}
	}
}

func (m *Manager) SetNotificationHandler(handler func(config.TunnelStatus)) {
	m.NotificationHandler = handler
}

func (m *Manager) callNotificationHandler(status config.TunnelStatus) {
	if m.NotificationHandler != nil {
		m.NotificationHandler(status)
	}
}

func (m *Manager) SetStatusChannel(ch chan config.TunnelStatus) {
	m.StatusChannel = ch
}

// SetProviderExpiration records how long each provider keeps a tunnel alive, in
// minutes, with 0 meaning "does not expire".
//
// Without this the ExpiresAt field is never populated, which is why the web
// dashboard's expiry warnings never fired: the frontend has had the countdown
// logic all along, but nothing ever sent it a deadline.
func (m *Manager) SetProviderExpiration(minutes map[string]int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.providerExpiration = make(map[string]int, len(minutes))
	for provider, mins := range minutes {
		m.providerExpiration[provider] = mins
	}
}

// expiresAtLocked returns the deadline for a tunnel started at startedAt, in
// Unix milliseconds, or 0 if the provider does not expire.
func (m *Manager) expiresAtLocked(provider config.Provider, startedAt time.Time) int64 {
	mins, ok := m.providerExpiration[string(provider)]
	if !ok || mins <= 0 {
		return 0
	}
	return startedAt.Add(time.Duration(mins) * time.Minute).UnixMilli()
}

// callStatusUpdate publishes a status change to whoever is listening.
//
// Every caller holds m.mu, so this must not block. The channel is buffered, but
// a blocking send would freeze the entire Manager behind the mutex if the
// consumer ever stalled or was never started -- taking the TUI down with it,
// since it shares the same lock. Dropping is the safe failure mode here: each
// message carries the full status, and clients resync over /api/tunnels.
func (m *Manager) callStatusUpdate(status config.TunnelStatus) {
	if m.StatusChannel == nil {
		return
	}

	select {
	case m.StatusChannel <- status:
	default:
		log.Printf("process: status channel full, dropped %q update for tunnel %s", status.State, status.ID)
	}
}

func (m *Manager) SetExpirationCallbacks(start func(string, string, string, time.Time), stop func(string)) {
	m.ExpirationCallbacks.OnStart = start
	m.ExpirationCallbacks.OnStop = stop
}

func (m *Manager) callExpirationStart(tunnelID, name, provider string, startedAt time.Time) {
	if m.ExpirationCallbacks.OnStart != nil {
		m.ExpirationCallbacks.OnStart(tunnelID, name, provider, startedAt)
	}
}

func (m *Manager) callExpirationStop(tunnelID string) {
	if m.ExpirationCallbacks.OnStop != nil {
		m.ExpirationCallbacks.OnStop(tunnelID)
	}
}

func NewManager() *Manager {
	return &Manager{
		processes: make(map[string]*ManagedProcess),
		providers: map[config.Provider]providers.Provider{
			config.ProviderCloudflared:  cloudflared.New(),
			config.ProviderTunnelmole:   tunnelmole.New(),
			config.ProviderLocalhostRun: ssh.NewLocalhostRun(),
			config.ProviderServeo:       ssh.NewServeo(),
			config.ProviderPinggy:       pinggy.New(),
			config.ProviderBore:         bore.New(),
		},
	}
}

func (m *Manager) CheckInstallation(providerType config.Provider) (needsInstall bool, autoInstall bool) {
	provider, ok := m.providers[providerType]
	if !ok {
		return false, false
	}

	installer, ok := provider.(providers.AutoInstaller)
	if !ok {
		return false, false
	}

	if installer.IsInstalled() {
		return false, true
	}

	return true, true
}

func (m *Manager) InstallProvider(providerType config.Provider) error {
	provider, ok := m.providers[providerType]
	if !ok {
		return fmt.Errorf("unknown provider: %s", providerType)
	}

	installer, ok := provider.(providers.AutoInstaller)
	if !ok {
		return fmt.Errorf("provider %s does not support auto-install", providerType)
	}

	return installer.Install(m.DownloadProgress)
}

func (m *Manager) Start(tunnel config.TunnelConfig, onUpdate func(config.TunnelStatus)) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if existing, ok := m.processes[tunnel.ID]; ok && existing.Process != nil {
		return fmt.Errorf("tunnel %s is already running", tunnel.ID)
	}

	provider, ok := m.providers[tunnel.Provider]
	if !ok {
		return fmt.Errorf("unknown provider: %s", tunnel.Provider)
	}

	if installer, ok := provider.(providers.AutoInstaller); ok && !installer.IsInstalled() {
		return fmt.Errorf("installing")
	}

	logBuffer := NewLogBuffer()
	mp := &ManagedProcess{
		Config:         tunnel,
		Provider:       provider,
		LogBuffer:      logBuffer,
		OnUpdate:       onUpdate,
		Status:         tunnel.Status(),
		logSubscribers: make(map[chan string]struct{}),
	}
	logBuffer.OnNewLine = func(line string) {
		mp.publishLog(line)
	}

	urlCapture := newURLCapture(provider, func(url string) { m.updateURL(tunnel.ID, url) })
	writer := io.MultiWriter(logBuffer, urlCapture)

	ctx := context.Background()
	proc, err := provider.Start(ctx, tunnel, writer)
	if err != nil {
		return err
	}
	startedAt := time.Now()
	mp.Process = proc
	mp.Status.State = config.TunnelStateStarting
	mp.Status.ExpiresAt = m.expiresAtLocked(tunnel.Provider, startedAt)

	m.processes[tunnel.ID] = mp

	if onUpdate != nil {
		onUpdate(mp.Status)
	}
	m.callStatusUpdate(mp.Status)

	go m.startupTimeoutMonitor(tunnel.ID, proc)
	go m.watchExit(tunnel.ID, proc)
	m.callExpirationStart(tunnel.ID, tunnel.Name, string(tunnel.Provider), startedAt)

	return nil
}

// publishLocked notifies listeners of mp's current status.
//
// The caller must hold m.mu. That is safe because callStatusUpdate never
// blocks; see its comment.
func (m *Manager) publishLocked(mp *ManagedProcess) {
	if mp.OnUpdate != nil {
		mp.OnUpdate(mp.Status)
	}
	m.callStatusUpdate(mp.Status)
}

// watchExit reports the outcome when a provider exits on its own.
//
// Nothing used to notice this at all: if cloudflared crashed, the UI kept
// showing the tunnel as online indefinitely and the child was never reaped.
//
// A deliberate Stop removes the tunnel from the map before tearing the process
// down, so this stays quiet for those and only speaks up for real failures.
func (m *Manager) watchExit(tunnelID string, proc *providers.Process) {
	<-proc.Exited()

	m.mu.Lock()
	mp, ok := m.processes[tunnelID]
	if !ok || mp.Process != proc {
		m.mu.Unlock()
		return
	}

	delete(m.processes, tunnelID)
	mp.closeLogSubscribers()

	mp.Status.PublicURL = ""
	if err := proc.Err(); err != nil {
		mp.Status.State = config.TunnelStateError
		mp.Status.ErrorMessage = fmt.Sprintf("tunnel process exited unexpectedly: %v", err)
	} else {
		mp.Status.State = config.TunnelStateStopped
		mp.Status.ErrorMessage = ""
	}

	status := mp.Status
	m.publishLocked(mp)
	m.mu.Unlock()

	m.callNotificationHandler(status)
	m.callExpirationStop(tunnelID)
}

const (
	// connectingAfter is how long a tunnel may sit in "starting" before it is
	// shown as "connecting", and giveUpAfter how much longer it gets to
	// produce a URL before being declared timed out.
	connectingAfter = 5 * time.Second
	giveUpAfter     = 25 * time.Second
)

// startupTimeoutMonitor gives up on a tunnel that never produces a URL.
//
// Every step re-checks that the process it was started for is still the one
// registered. Previously this just slept and then mutated whatever it found:
// stopping a tunnel and starting it again inside the 30s window let the stale
// monitor mark the *new* process as timed out.
func (m *Manager) startupTimeoutMonitor(tunnelID string, proc *providers.Process) {
	select {
	case <-proc.Exited():
		return // watchExit reports this
	case <-time.After(connectingAfter):
	}

	m.mu.Lock()
	mp, ok := m.processes[tunnelID]
	if !ok || mp.Process != proc {
		m.mu.Unlock()
		return
	}
	if mp.Status.PublicURL == "" && mp.Status.State != config.TunnelStateOnline {
		mp.Status.State = config.TunnelStateConnecting
		m.publishLocked(mp)
	}
	m.mu.Unlock()

	select {
	case <-proc.Exited():
		return
	case <-time.After(giveUpAfter):
	}

	m.mu.Lock()
	mp, ok = m.processes[tunnelID]
	if !ok || mp.Process != proc {
		m.mu.Unlock()
		return
	}
	if mp.Status.State != config.TunnelStateConnecting && mp.Status.State != config.TunnelStateStarting {
		m.mu.Unlock()
		return
	}

	// Removed before teardown so watchExit stays quiet about this exit.
	delete(m.processes, tunnelID)
	mp.closeLogSubscribers()

	mp.Status.State = config.TunnelStateTimeout
	mp.Status.ErrorMessage = "Connection timed out after 30 seconds"
	status := mp.Status
	m.publishLocked(mp)
	m.mu.Unlock()

	m.callNotificationHandler(status)
	m.callExpirationStop(tunnelID)

	proc.Stop()
}

// Stop terminates a tunnel and waits for its process to actually be gone.
//
// The tunnel is removed from the map up front so watchExit treats the exit as
// deliberate, and the teardown itself happens with m.mu released: Process.Stop
// blocks until the process dies, and holding the lock across that would stall
// every other Manager operation, TUI included.
func (m *Manager) Stop(tunnelID string) error {
	m.mu.Lock()

	mp, ok := m.processes[tunnelID]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("tunnel %s is not running", tunnelID)
	}

	delete(m.processes, tunnelID)
	mp.closeLogSubscribers()

	proc := mp.Process
	mp.Status.State = config.TunnelStateStopping
	mp.Status.PublicURL = ""
	stopping := mp.Status
	m.publishLocked(mp)
	m.mu.Unlock()

	m.callNotificationHandler(stopping)

	if proc != nil {
		proc.Stop()
	}

	m.mu.Lock()
	mp.Status.State = config.TunnelStateStopped
	mp.Status.ErrorMessage = ""
	m.publishLocked(mp)
	m.mu.Unlock()

	m.callExpirationStop(tunnelID)

	return nil
}

func (m *Manager) GetStatus(tunnelID string) (config.TunnelStatus, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	mp, ok := m.processes[tunnelID]
	if !ok {
		return config.TunnelStatus{}, false
	}
	return mp.Status, true
}

func (m *Manager) updateURL(tunnelID, url string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	mp, ok := m.processes[tunnelID]
	if !ok {
		return
	}

	mp.PublicURL = url
	mp.Status.PublicURL = url
	mp.Status.State = config.TunnelStateOnline

	if mp.OnUpdate != nil {
		mp.OnUpdate(mp.Status)
	}

	m.callNotificationHandler(mp.Status)
	m.callStatusUpdate(mp.Status)
}

func (m *Manager) GetLogs(tunnelID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	mp, ok := m.processes[tunnelID]
	if !ok {
		return []string{}
	}

	return mp.LogBuffer.GetLines()
}

func (m *Manager) SubscribeLogs(tunnelID string) (<-chan string, func()) {
	m.mu.RLock()
	mp, ok := m.processes[tunnelID]
	m.mu.RUnlock()
	if !ok {
		return nil, func() {}
	}

	ch, cancel := mp.addLogSubscriber()
	return ch, cancel
}

func (m *Manager) IsRunning(tunnelID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	mp, ok := m.processes[tunnelID]
	return ok && (mp.Status.State == config.TunnelStateOnline || mp.Status.State == config.TunnelStateConnecting || mp.Status.State == config.TunnelStateStarting)
}

// StopAll terminates every tunnel and waits for them all to be gone.
//
// Waiting is the point: this runs on shutdown, and returning early left the
// provider processes running after ftm exited, with the user's world still
// exposed by a tunnel they thought they had closed.
func (m *Manager) StopAll() {
	m.mu.Lock()
	procs := make([]*providers.Process, 0, len(m.processes))
	for _, mp := range m.processes {
		mp.closeLogSubscribers()
		if mp.Process != nil {
			procs = append(procs, mp.Process)
		}
	}
	m.processes = make(map[string]*ManagedProcess)
	m.mu.Unlock()

	// In parallel: each Stop waits out its own process, and on shutdown those
	// waits should not be serialised.
	var wg sync.WaitGroup
	for _, proc := range procs {
		wg.Add(1)
		go func(p *providers.Process) {
			defer wg.Done()
			p.Stop()
		}(proc)
	}
	wg.Wait()
}
