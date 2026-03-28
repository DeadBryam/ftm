package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/sthbryan/ftm/internal/clipboard"
	"github.com/sthbryan/ftm/internal/config"
)

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.Keys.Quit):
		return m.handleQuit()

	case key.Matches(msg, m.Keys.Back):
		return m.handleBack()

	case key.Matches(msg, m.Keys.Help):
		m.Help.ShowAll = !m.Help.ShowAll
		return m, nil
	}

	switch m.State {
	case viewList:
		return m.handleListKey(msg)
	case viewLogs:
		return m.handleLogsKey(msg)
	case viewAddForm, viewEditForm:
		return m.handleFormKey(msg)
	case viewDownloading:
		return m.handleDownloadingKey(msg)
	}

	return m, nil
}

func (m *Model) handleQuit() (tea.Model, tea.Cmd) {
	if m.State == viewList {
		return m, tea.Quit
	}
	m.State = viewList
	return m, nil
}

func (m *Model) handleBack() (tea.Model, tea.Cmd) {
	if m.State != viewList {
		m.State = viewList
		m.editingTunnelID = ""
	}
	return m, nil
}

func (m *Model) handleDownloadingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.Keys.Back) || key.Matches(msg, m.Keys.Quit) {
		m.State = viewList
	}
	return m, nil
}

func (m *Model) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.Keys.Up):
		m.moveCursorUp()

	case key.Matches(msg, m.Keys.Down):
		m.moveCursorDown()

	case key.Matches(msg, m.Keys.Enter), key.Matches(msg, m.Keys.Toggle):
		return m.handleListToggle()

	case key.Matches(msg, m.Keys.Logs):
		return m.handleListLogs()

	case key.Matches(msg, m.Keys.Copy):
		m.handleListCopy()

	case key.Matches(msg, m.Keys.Web):
		m.openDashboard()

	case key.Matches(msg, m.Keys.Config):
		m.openConfigDir()

	case key.Matches(msg, m.Keys.Add):
		m.startAddForm()

	case key.Matches(msg, m.Keys.Edit):
		return m.startEditForm()

	case key.Matches(msg, m.Keys.Delete):
		return m.handleListDelete()
	}

	return m, nil
}

func (m *Model) moveCursorUp() {
	if m.Cursor > 0 {
		m.Cursor--
	}
}

func (m *Model) moveCursorDown() {
	if m.Cursor < len(m.Items)-1 {
		m.Cursor++
	}
}

func (m *Model) handleListToggle() (tea.Model, tea.Cmd) {
	if item, ok := m.selectedItem(); ok {
		if m.App.Manager.IsRunning(item.Tunnel.ID) {
			m.playBeep()
			return m, m.stopTunnel(item)
		}
		m.playBeep()
		return m, m.startTunnel(item)
	}
	return m, nil
}

func (m *Model) handleListLogs() (tea.Model, tea.Cmd) {
	if item, ok := m.selectedItem(); ok {
		m.SelectedTunnel = item.Tunnel.ID
		m.State = viewLogs
		m.updateLogViewport()
	}
	return m, nil
}

func (m *Model) handleListCopy() {
	if item, ok := m.selectedItem(); ok {
		m.copyTunnelURL(item)
	}
}

func (m *Model) startAddForm() {
	m.State = viewAddForm
	m.FormFocus = 0
	m.FormValues = FormData{
		Provider: string(config.ProviderCloudflared),
		Port:     "30000",
	}
}

func (m *Model) startEditForm() (tea.Model, tea.Cmd) {
	if item, ok := m.selectedItem(); ok {
		if item.Status.State != config.TunnelStateStopped {
			m.showMessage("Stop tunnel first to edit")
			return m, nil
		}
		m.State = viewEditForm
		m.editingTunnelID = item.Tunnel.ID
		m.FormFocus = 0
		m.FormValues = FormData{
			ID:       item.Tunnel.ID,
			Name:     item.Tunnel.Name,
			Provider: string(item.Tunnel.Provider),
			Port:     fmt.Sprintf("%d", item.Tunnel.LocalPort),
		}
	}
	return m, nil
}

func (m *Model) handleListDelete() (tea.Model, tea.Cmd) {
	if item, ok := m.selectedItem(); ok {
		m.App.Manager.Stop(item.Tunnel.ID)
		m.App.Config.RemoveTunnel(item.Tunnel.ID)
		m.App.SaveConfig()
		m.refreshItems()
		if m.Cursor >= len(m.Items) && m.Cursor > 0 {
			m.Cursor--
		}
		m.showMessage("Tunnel deleted")
	}
	return m, nil
}

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

func (m *Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.State != viewList {
		return m, nil
	}

	itemHeight := 3
	headerHeight := 4
	clickedIdx := (msg.Y - headerHeight) / itemHeight

	switch msg.Type {
	case tea.MouseLeft:
		if clickedIdx >= 0 && clickedIdx < len(m.Items) {
			m.Cursor = clickedIdx
		}

	case tea.MouseWheelUp:
		if m.Cursor > 0 {
			m.Cursor--
		}

	case tea.MouseWheelDown:
		if m.Cursor < len(m.Items)-1 {
			m.Cursor++
		}
	}

	return m, nil
}

func (m *Model) copyTunnelURL(item TunnelItem) {
	if item.Status.PublicURL != "" {
		clipboard.Write(item.Status.PublicURL)
		m.showMessage("Copied URL!")
		return
	}

	m.showMessage("No URL available - start tunnel first")
}

func (m *Model) openDashboard() {
	if err := m.App.OpenDashboard(); err != nil {
		m.showMessage("Error opening dashboard: " + err.Error())
		return
	}
	m.showMessage("Dashboard opened in browser")
}

func (m *Model) openConfigDir() {
	if err := m.App.OpenConfigDir(); err != nil {
		m.showMessage("Error opening config folder: " + err.Error())
		return
	}
	m.showMessage("Config folder opened")
}

func (m *Model) playBeep() {
	fmt.Fprint(os.Stdout, Bell)
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},
		{k.Toggle, k.Logs, k.Copy, k.Web},
		{k.Add, k.Delete, k.Config},
		{k.Back, k.Help, k.Quit},
	}
}
