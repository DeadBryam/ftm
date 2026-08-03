package views

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/sthbryan/ftm/internal/app/ui"
	"github.com/sthbryan/ftm/internal/i18n"
)

type FormView struct {
	Width      int
	Height     int
	Focus      int
	IsEditMode bool
	Name       string
	Provider   string
	Port       string
}

func NewFormView() *FormView {
	return &FormView{}
}

const (
	formWidth   = 54
	submitFocus = 4
)

func (f *FormView) Render() string {
	t := ui.ThemeDefault
	inner := ui.PanelInner(formWidth)

	title, subtitle := i18n.T("new_tunnel"), i18n.T("new_tunnel_desc")
	if f.IsEditMode {
		title, subtitle = i18n.T("edit_tunnel"), i18n.T("edit_tunnel_desc")
	}

	dim := lipgloss.NewStyle().Foreground(t.TextDim).Width(inner)

	body := strings.Join([]string{
		dim.Render(subtitle),
		"",
		f.field(inner, 0, i18n.T("name_label"), i18n.T("tunnel_name_hint"), i18n.T("type_hint"), f.Name),
		"",
		f.field(inner, 1, i18n.T("provider_label"), i18n.T("provider_hint"), i18n.T("arrow_hint"), f.Provider),
		"",
		f.field(inner, 2, i18n.T("local_port"), i18n.T("port_hint"), i18n.T("numbers_hint"), f.Port),
		"",
		lipgloss.NewStyle().Width(inner).Align(lipgloss.Center).Render(f.submitButton()),
		"",
		dim.Align(lipgloss.Center).Render(i18n.T("form_nav_hint")),
	}, "\n")

	panel := ui.Panel(title, body, formWidth, lipgloss.Height(body)+ui.PanelChrome, t.Gold)

	return ui.Overlay("", panel, f.Width, f.Height)
}

func (f *FormView) field(inner, index int, label, hint, focusedHint, value string) string {
	t := ui.ThemeDefault
	focused := f.Focus == index

	if value == "" {
		value = "…"
	}

	labelStyle := lipgloss.NewStyle().Foreground(t.TextDim)
	border, foreground := t.Bronze, t.TextDim

	if focused {
		label = "▸ " + label
		hint = focusedHint
		labelStyle = labelStyle.Foreground(t.Gold).Bold(true)
		border, foreground = t.Gold, t.Text

		if index == 1 {
			value = "‹ " + value + " ›"
		}
	}

	input := lipgloss.NewStyle().
		Width(inner).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Foreground(foreground).
		Render(value)

	return strings.Join([]string{
		labelStyle.Render(label),
		input,
		lipgloss.NewStyle().Foreground(t.Bronze).Render(hint),
	}, "\n")
}

func (f *FormView) submitButton() string {
	t := ui.ThemeDefault

	label := i18n.T("submit_new")
	if f.IsEditMode {
		label = i18n.T("submit_edit")
	}

	background, foreground := t.Button, t.ButtonText
	if f.Focus == submitFocus {
		background, foreground = t.ButtonActive, t.ButtonActiveText
	}

	return lipgloss.NewStyle().
		Background(background).
		Foreground(foreground).
		Bold(true).
		Padding(0, 4).
		Render(label)
}
