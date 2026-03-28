package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

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
