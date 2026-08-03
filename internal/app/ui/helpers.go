package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const HeaderMargin = 4

func Repeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(s, count)
}

func Clamp(v, min int) int {
	if v < min {
		return min
	}
	return v
}

func Center(s string, width int) string {
	pad := (width - lipgloss.Width(s)) / 2
	if pad < 0 {
		pad = 0
	}
	return strings.Repeat(" ", pad) + s
}
