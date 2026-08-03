package views

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/sthbryan/ftm/internal/app/ui"
	"github.com/sthbryan/ftm/internal/i18n"
)

type LogsView struct {
	Width      int
	Height     int
	TunnelName string
	Content    string
}

func NewLogsView() *LogsView {
	return &LogsView{}
}

func (l *LogsView) Render() string {
	t := ui.ThemeDefault

	footer := lipgloss.NewStyle().
		Foreground(t.TextDim).
		Render(i18n.T("logs_nav_hint"))

	bodyHeight := ui.Clamp(l.Height-lipgloss.Height(footer)-1, minBodyHeight)

	content := lipgloss.NewStyle().
		Foreground(t.Text).
		Render(ui.TailLines(strings.TrimRight(l.Content, "\n"), bodyHeight-ui.PanelChrome))

	title := i18n.T("tunnel_logs")
	if l.TunnelName != "" {
		title += "  ·  " + l.TunnelName
	}

	return ui.Panel(title, content, l.Width, bodyHeight, t.Gold) + "\n" + footer
}
