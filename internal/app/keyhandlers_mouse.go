package app

import (
	tea "charm.land/bubbletea/v2"
)

func (m *Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.State != viewList {
		return m, nil
	}

	mouse := msg.Mouse()
	itemHeight := 3
	headerHeight := 4
	clickedIdx := (mouse.Y - headerHeight) / itemHeight

	switch msg := msg.(type) {
	case tea.MouseClickMsg:
		if msg.Button == tea.MouseLeft {
			if clickedIdx >= 0 && clickedIdx < len(m.Items) {
				m.Cursor = clickedIdx
			}
		}

	case tea.MouseWheelMsg:
		switch msg.Button {
		case tea.MouseWheelUp:
			if m.Cursor > 0 {
				m.Cursor--
			}
		case tea.MouseWheelDown:
			if m.Cursor < len(m.Items)-1 {
				m.Cursor++
			}
		}
	}

	return m, nil
}
