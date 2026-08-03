package ui

import (
	"charm.land/lipgloss/v2"
)

func Overlay(background, foreground string, width, height int) string {
	if width <= 0 || height <= 0 {
		return foreground
	}

	x := Clamp((width-lipgloss.Width(foreground))/2, 0)
	y := Clamp((height-lipgloss.Height(foreground))/2, 0)

	return lipgloss.NewCanvas(width, height).
		Compose(lipgloss.NewCompositor(
			lipgloss.NewLayer(background),
			lipgloss.NewLayer(foreground).X(x).Y(y).Z(1),
		)).
		Render()
}
