package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

const HeaderMargin = 4

// Repeat is strings.Repeat with a floor of zero.
//
// Layout widths here are computed by subtraction (total minus the rendered
// pieces), which goes negative on a narrow terminal. strings.Repeat panics on a
// negative count, so every layout call site must go through this.
func Repeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(s, count)
}

// Clamp returns v bounded below by min.
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
