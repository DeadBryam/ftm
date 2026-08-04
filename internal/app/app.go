package app

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"

	tea "charm.land/bubbletea/v2"

	"github.com/sthbryan/ftm/internal/config"
	"github.com/sthbryan/ftm/internal/i18n"
	"github.com/sthbryan/ftm/internal/notifications"
	"github.com/sthbryan/ftm/internal/process"
	"github.com/sthbryan/ftm/internal/providers"
	"github.com/sthbryan/ftm/internal/web"
)

type App struct {
	Config            *config.Config
	Manager           *process.Manager
	WebServer         *web.Server
	DownloadProgress  chan providers.DownloadProgress
	ExpirationMonitor *notifications.ExpirationMonitor

	shutdownOnce sync.Once
}

func New() (*App, error) {

	if err := i18n.Load(); err != nil {
		return nil, fmt.Errorf("failed to load translations: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	i18n.InitFromConfig(cfg)

	manager := process.NewManager()

	app := &App{
		Config:           cfg,
		Manager:          manager,
		DownloadProgress: make(chan providers.DownloadProgress, 10),
	}

	app.Manager.SetProgressChannel(app.DownloadProgress)
	app.Manager.SetProviderExpiration(cfg.ProviderExpirationMinutes)
	app.Manager.SetVisitorTracking(cfg.CountVisitors)
	notifications.Init()
	notifications.SetSoundEnabled(cfg.NotificationSound)
	notifications.SetNotificationsEnabled(cfg.NotificationsStatus == config.NotificationGranted)

	expConfig := notifications.ExpirationConfig{
		Thresholds:                cfg.ExpirationThresholds,
		ProviderExpirationMinutes: cfg.ProviderExpirationMinutes,
	}
	app.ExpirationMonitor = notifications.NewExpirationMonitor(expConfig, func(name string, mins int) {
		if !app.shouldUseNativeNotifications() {
			return
		}
		if mins == 0 {
			notifications.NotifyTunnelExpired(name)
		} else {
			notifications.NotifyTunnelExpiring(name, mins)
		}
	})

	app.Manager.SetNotificationHandler(func(status config.TunnelStatus) {
		if !app.shouldUseNativeNotifications() {
			return
		}
		switch status.State {
		case config.TunnelStateOnline:
			notifications.NotifyTunnelOnline(status.Name, status.PublicURL)
		case config.TunnelStateError:
			notifications.NotifyTunnelError(status.Name, status.ErrorMessage)
		case config.TunnelStateTimeout:
			notifications.NotifyTunnelTimeout(status.Name)
		case config.TunnelStateStopping:
			notifications.NotifyTunnelStopped(status.Name)
		}
	})

	app.Manager.SetVisitorHandler(func(status config.TunnelStatus) {
		if app.WebServer != nil {
			app.WebServer.BroadcastNewVisitor(status)
		}
		if !app.shouldUseNativeNotifications() {
			return
		}
		notifications.NotifyNewVisitor(status.Name, status.Visitors)
	})

	app.Manager.SetExpirationCallbacks(
		app.ExpirationMonitor.Start,
		func(id string) { app.ExpirationMonitor.Stop(id) },
	)

	return app, nil
}

func (a *App) Run() error {
	if len(a.Config.Tunnels) == 0 {
		a.createDefaultTunnels()
	}

	if file := redirectLog(); file != nil {
		defer file.Close()
	}

	if err := a.StartWebServer(); err != nil {
		return fmt.Errorf("failed to start web server: %w", err)
	}

	model := NewModel(a)
	p := tea.NewProgram(model)

	_, err := p.Run()
	return err
}

func (a *App) StartWebServer() error {
	if a.WebServer != nil {
		return nil
	}

	a.WebServer = web.NewServer(a.Manager, a.Config)
	a.Manager.SetStatusChannel(a.WebServer.StatusChannel)

	if err := a.WebServer.Start(); err != nil {
		a.WebServer = nil
		return err
	}

	return nil
}

func (a *App) OpenDashboard() error {
	if a.WebServer == nil {
		return fmt.Errorf("web server not started")
	}

	url := a.WebServer.URL()

	switch runtime.GOOS {
	case "darwin":
		return startAndReap(exec.Command("open", url))
	case "windows":
		return startAndReap(exec.Command("rundll32", "url.dll,FileProtocolHandler", url))
	default:
		return startAndReap(exec.Command("xdg-open", url))
	}
}

func (a *App) OpenConfigDir() error {
	path := config.ConfigDir()

	switch runtime.GOOS {
	case "darwin":
		return startAndReap(exec.Command("open", path))
	case "windows":
		return startAndReap(exec.Command("explorer", path))
	default:
		return startAndReap(exec.Command("xdg-open", path))
	}
}

var onReaped = func(*exec.Cmd) {}

func startAndReap(cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		_ = cmd.Wait()
		onReaped(cmd)
	}()

	return nil
}

func (a *App) createDefaultTunnels() {
	a.Config.Tunnels = []config.TunnelConfig{
		{
			ID:        "foundry-default",
			Name:      "Foundry VTT (Default)",
			Provider:  config.ProviderCloudflared,
			LocalPort: 30000,
		},
	}
	a.Config.Save()
}

func (a *App) SaveConfig() error {
	notifications.SetSoundEnabled(a.Config.NotificationSound)
	notifications.SetNotificationsEnabled(a.Config.NotificationsStatus == config.NotificationGranted)
	return a.Config.Save()
}

func (a *App) shouldUseNativeNotifications() bool {
	if a.WebServer == nil {
		return true
	}
	return a.WebServer.ClientCount() == 0
}

func (a *App) Shutdown() {
	a.shutdownOnce.Do(func() {

		if a.Manager != nil {
			a.Manager.StopAll()
		}
		if a.WebServer != nil {
			if err := a.WebServer.Stop(); err != nil {
				fmt.Fprintf(os.Stderr, "Error stopping web server: %v\n", err)
			}
		}
		if a.ExpirationMonitor != nil {
			a.ExpirationMonitor.StopAll()
		}
	})
}
