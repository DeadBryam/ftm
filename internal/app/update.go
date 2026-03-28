package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/sthbryan/ftm/internal/providers"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.LogViewport.Width = msg.Width - 4
		m.LogViewport.Height = msg.Height - 8
		return m, nil

	case tickMsg:
		return m.handleTick()

	case downloadProgressMsg:
		return m.handleDownloadProgress(msg)

	case statusUpdateMsg:
		return m.handleStatusUpdate(msg)
	}

	return m, nil
}

func (m *Model) handleTick() (tea.Model, tea.Cmd) {
	m.refreshItems()

	if m.MessageTimer > 0 {
		m.MessageTimer--
		if m.MessageTimer == 0 {
			m.Message = ""
		}
	}
	return m, tickCmd()
}

func (m *Model) handleDownloadProgress(msg downloadProgressMsg) (tea.Model, tea.Cmd) {
	m.DownloadProgress = providers.DownloadProgress(msg)
	if m.DownloadProgress.Done {
		m.State = viewList
		if m.PendingTunnel != nil {
			for _, item := range m.Items {
				if ti, ok := item.(TunnelItem); ok && ti.Tunnel.ID == m.PendingTunnel.ID {
					m.PendingTunnel = nil
					m.showMessage("Install complete! Starting tunnel...")
					return m, m.startTunnel(ti)
				}
			}
			m.PendingTunnel = nil
		}
		m.showMessage("Download complete!")
	}
	return m, m.checkDownloadProgress()
}

func (m *Model) handleStatusUpdate(msg statusUpdateMsg) (tea.Model, tea.Cmd) {
	m.refreshItems()
	if msg.status.ErrorMessage != "" {
		m.playBeep()
		m.showMessage("Error: " + msg.status.ErrorMessage)
	}
	return m, nil
}
