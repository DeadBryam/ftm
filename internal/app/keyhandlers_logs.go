package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleLogsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.LogViewport, cmd = m.LogViewport.Update(msg)
	m.updateLogViewport()
	return m, cmd
}

func (m *Model) updateLogViewport() {
	if m.SelectedTunnel == "" {
		return
	}

	logs := m.App.Manager.GetLogs(m.SelectedTunnel)
	content := strings.Join(logs, "\n")
	m.LogViewport.SetContent(content)
	m.LogViewport.GotoBottom()
}
