package app

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// entersText reports whether the current view has a focused text field, in
// which case letters have to reach it rather than being read as shortcuts.
func (m *Model) entersText() bool {
	return m.State == viewAddForm || m.State == viewEditForm
}

func (m *Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	// Typing a tunnel name must not quit the app, so while a field has focus
	// only Ctrl+C and Esc are handled globally.
	if m.entersText() {
		switch {
		case key.Matches(msg, m.Keys.ForceQuit):
			return m, tea.Quit
		case key.Matches(msg, m.Keys.Back):
			return m.handleBack()
		}
		return m.handleFormKey(msg)
	}

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
	case viewDownloading:
		return m.handleDownloadingKey(msg)
	case viewSettings:
		return m.handleSettingsKey(msg)
	}

	return m, nil
}

func (m *Model) handleQuit() (tea.Model, tea.Cmd) {
	if m.State == viewSettings {
		m.saveSettings()
	}
	return m, tea.Quit
}

func (m *Model) handleBack() (tea.Model, tea.Cmd) {
	if m.State == viewSettings {
		m.saveSettings()
	}
	if m.State != viewList {
		m.State = viewList
		m.editingTunnelID = ""
	}
	return m, nil
}

func (m *Model) handleDownloadingKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.Keys.Back) || key.Matches(msg, m.Keys.Quit) {
		m.State = viewList
	}
	return m, nil
}
