package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"foundry-tunnel/internal/config"
	"foundry-tunnel/internal/process"
	"foundry-tunnel/internal/providers"
)

type App struct {
	Config           *config.Config
	Manager          *process.Manager
	DownloadProgress chan providers.DownloadProgress
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	app := &App{
		Config:           cfg,
		Manager:          process.NewManager(),
		DownloadProgress: make(chan providers.DownloadProgress, 10),
	}

	app.Manager.SetProgressChannel(app.DownloadProgress)

	return app, nil
}

func (a *App) Run() error {
	if len(a.Config.Tunnels) == 0 {
		a.createDefaultTunnels()
	}

	model := NewModel(a)
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}

func (a *App) createDefaultTunnels() {
	a.Config.Tunnels = []config.TunnelConfig{
		{
			ID:        "foundry-default",
			Name:      "Foundry VTT (Default)",
			Provider:  config.ProviderPlayitgg,
			LocalPort: 30000,
			AutoStart: false,
		},
	}
	a.Config.Save()
}

func (a *App) SaveConfig() error {
	return a.Config.Save()
}
