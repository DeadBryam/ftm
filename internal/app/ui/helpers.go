package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const HeaderMargin = 4

func Center(s string, width int) string {
	pad := (width - lipgloss.Width(s)) / 2
	if pad < 0 {
		pad = 0
	}
	return strings.Repeat(" ", pad) + s
}
