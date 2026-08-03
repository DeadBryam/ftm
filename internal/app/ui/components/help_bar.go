package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sthbryan/ftm/internal/app/ui"
)

// Shortcut is one entry in the help bar.
type Shortcut struct {
	Keys  string
	Label string
}

// HelpBar renders the shortcut hints under the tunnel list.
//
// Entries are supplied by the caller rather than hardcoded here: they come from
// the KeyMap, so a rebound key updates the hint automatically instead of
// leaving the bar advertising a shortcut that no longer works.
type HelpBar struct {
	Shortcuts []Shortcut
	Width     int
}

func NewHelpBar() *HelpBar {
	return &HelpBar{}
}

const helpSeparator = "  •  "

func (h *HelpBar) Render() string {
	if len(h.Shortcuts) == 0 {
		return ""
	}

	keyStyle := lipgloss.NewStyle().Foreground(ui.ThemeDefault.Gold)
	labelStyle := lipgloss.NewStyle().Foreground(ui.ThemeDefault.TextDim)
	separator := labelStyle.Render(helpSeparator)

	entries := make([]string, 0, len(h.Shortcuts))
	for _, s := range h.Shortcuts {
		entries = append(entries, keyStyle.Render(s.Keys)+" "+labelStyle.Render(s.Label))
	}

	// Wrapped to the available width rather than split down the middle, so a
	// narrow terminal does not push half the bar off-screen.
	var b strings.Builder
	lineWidth := 0

	for i, entry := range entries {
		entryWidth := lipgloss.Width(entry)

		switch {
		case i == 0:
			lineWidth = entryWidth
		case h.Width > 0 && lineWidth+lipgloss.Width(separator)+entryWidth > h.Width:
			b.WriteString("\n")
			lineWidth = entryWidth
		default:
			b.WriteString(separator)
			lineWidth += lipgloss.Width(separator) + entryWidth
		}

		b.WriteString(entry)
	}

	return b.String()
}
