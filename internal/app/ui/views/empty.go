package views

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/sthbryan/ftm/internal/app/ui"
	"github.com/sthbryan/ftm/internal/i18n"
)

type EmptyState struct {
	Height    int
	Width     int
	Dashboard string
	Sessions  int
}

func NewEmptyState() *EmptyState {
	return &EmptyState{}
}

const emptyWidth = 62

func (e *EmptyState) Render() string {
	t := ui.ThemeDefault
	inner := ui.PanelInner(emptyWidth)

	centered := lipgloss.NewStyle().Width(inner).Align(lipgloss.Center)

	body := strings.Join([]string{
		centered.Foreground(t.Text).Render(i18n.T("no_tunnels_yet")),
		centered.Foreground(t.TextDim).Render(i18n.T("tunnels_desc")),
		"",
		centered.Render(lipgloss.NewStyle().
			Background(t.ButtonActive).
			Foreground(t.ButtonActiveText).
			Bold(true).
			Padding(0, 3).
			Render(i18n.T("create_first"))),
		"",
		centered.Foreground(t.TextDim).Render(i18n.T("press_a_hint")),
		"",
		centered.Foreground(t.Bronze).Render(
			i18n.T("tip_dashboard") + " " + e.Dashboard + "  •  ws " + fmt.Sprintf("%d", e.Sessions)),
	}, "\n")

	panel := ui.Panel(i18n.T("welcome_title"), body, emptyWidth, lipgloss.Height(body)+ui.PanelChrome, t.Gold)

	return ui.Overlay("", panel, e.Width, e.Height)
}
