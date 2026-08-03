package ui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type Theme struct {
	Gold       color.Color
	Bronze     color.Color
	Text       color.Color
	TextDim    color.Color
	Online     color.Color
	Offline    color.Color
	Connecting color.Color
	Error      color.Color
	Danger     color.Color
	Stopped    color.Color
	Success    color.Color
}

func DefaultTheme() *Theme {
	return &Theme{
		Gold:       lipgloss.Color("#c9a227"),
		Bronze:     lipgloss.Color("#8b7355"),
		Text:       lipgloss.Color("#ffffff"),
		TextDim:    lipgloss.Color("#9a9590"),
		Online:     lipgloss.Color("#1e3a2f"),
		Offline:    lipgloss.Color("#2a2824"),
		Connecting: lipgloss.Color("#3a3020"),
		Error:      lipgloss.Color("#3a2020"),
		Danger:     lipgloss.Color("#ff6b6b"),
		Stopped:    lipgloss.Color("#3a3a3a"),
		Success:    lipgloss.Color("#7cb69d"),
	}
}

var ThemeDefault = DefaultTheme()
