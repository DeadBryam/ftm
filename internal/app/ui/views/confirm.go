package views

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/sthbryan/ftm/internal/app/ui"
	"github.com/sthbryan/ftm/internal/i18n"
)

type ConfirmView struct {
	Width   int
	Height  int
	Title   string
	Message string
	Danger  bool
}

func NewConfirmView() *ConfirmView {
	return &ConfirmView{}
}

func (c *ConfirmView) Render() string {
	t := ui.ThemeDefault

	accent := t.Gold
	if c.Danger {
		accent = t.Danger
	}

	title := lipgloss.NewStyle().
		Foreground(accent).
		Bold(true).
		Render(c.Title)

	message := lipgloss.NewStyle().
		Foreground(t.Text).
		Render(c.Message)

	confirm := lipgloss.NewStyle().Foreground(accent).Bold(true).Render(i18n.T("confirm_yes"))
	cancel := lipgloss.NewStyle().Foreground(t.TextDim).Render(i18n.T("confirm_no"))

	body := strings.Join([]string{title, "", message, "", confirm + "    " + cancel}, "\n")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accent).
		Padding(1, 3).
		Render(body)

	boxHeight := lipgloss.Height(box)
	top := ui.Clamp((c.Height-boxHeight)/2, 1)

	var b strings.Builder
	b.WriteString(ui.Repeat("\n", top))
	for _, line := range strings.Split(box, "\n") {
		b.WriteString(ui.Center(line, c.Width))
		b.WriteString("\n")
	}

	return b.String()
}
