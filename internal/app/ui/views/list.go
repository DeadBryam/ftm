package views

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/sthbryan/ftm/internal/app/ui"
	"github.com/sthbryan/ftm/internal/app/ui/components"
	"github.com/sthbryan/ftm/internal/i18n"
	"github.com/sthbryan/ftm/internal/version"
)

type TunnelViewData struct {
	Name        string
	Provider    string
	LocalPort   int
	StatusState int
	StatusMsg   string
	PublicURL   string
	ErrorMsg    string
}

type ListView struct {
	Width          int
	Height         int
	Items          []TunnelViewData
	Cursor         int
	Message        string
	Dashboard      string
	Sessions       int
	TwoColumnLimit int
	UpdateBadge    string
	Shortcuts      []components.Shortcut

	ListTop int

	VisibleRows int
	FirstItem   int
}

const ItemHeight = 1

const (
	singleColumnChrome = 6
	twoColumnChrome    = 7
)

func (l *ListView) rowsAvailable(chrome int) int {
	if l.Height <= 0 {
		return 0
	}

	rows := (l.Height - l.ListTop - chrome) / ItemHeight
	if rows < 1 {
		return 1
	}

	return rows
}

func NewListView() *ListView {
	return &ListView{
		TwoColumnLimit: 100,
	}
}

func (l *ListView) Render() string {
	if len(l.Items) == 0 {
		return ""
	}

	if l.Width >= l.TwoColumnLimit {
		return l.twoColumn()
	}

	return l.singleColumn()
}

func (l *ListView) twoColumn() string {
	var b strings.Builder

	gold := ui.ThemeDefault.Gold
	bronze := ui.ThemeDefault.Bronze
	text := ui.ThemeDefault.Text
	textDim := ui.ThemeDefault.TextDim

	title := lipgloss.NewStyle().
		Foreground(gold).
		Bold(true).
		Render(i18n.T("app_name_tui"))
	versionStr := lipgloss.NewStyle().
		Foreground(textDim).
		Render("v" + version.Version + "  ws:" + fmt.Sprintf("%d", l.Sessions))

	b.WriteString(title)
	b.WriteString(ui.Repeat(" ", l.Width-lipgloss.Width(title)-lipgloss.Width(versionStr)-ui.HeaderMargin))
	b.WriteString(versionStr)
	b.WriteString("\n")

	if l.UpdateBadge != "" {
		badgeStyle := lipgloss.NewStyle().Foreground(ui.ThemeDefault.Gold).Bold(true)
		b.WriteString(badgeStyle.Render(l.UpdateBadge))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	leftWidth := int(float64(l.Width) * 0.4)
	rightWidth := l.Width - leftWidth - 3

	leftHeader := lipgloss.NewStyle().
		Bold(true).
		Foreground(text).
		Render(i18n.T("connections"))

	rightHeader := lipgloss.NewStyle().
		Bold(true).
		Foreground(text).
		Render(i18n.T("selected_tunnel"))

	b.WriteString(leftHeader)
	b.WriteString(ui.Repeat(" ", leftWidth-lipgloss.Width(leftHeader)+3))
	b.WriteString(rightHeader)
	b.WriteString("\n")

	dividerStyle := lipgloss.NewStyle().Foreground(bronze)
	b.WriteString(dividerStyle.Render(ui.Repeat("─", l.Width-2)))
	b.WriteString("\n")

	l.ListTop = 4
	if l.UpdateBadge != "" {
		l.ListTop++
	}
	l.VisibleRows = l.rowsAvailable(twoColumnChrome)

	listContent := l.renderTunnelList(leftWidth - 2)
	detailContent := l.renderDetailPanel(rightWidth - 2)

	listLines := strings.Split(listContent, "\n")
	detailLines := strings.Split(detailContent, "\n")

	maxLines := len(listLines)
	if detail := len(detailLines); detail > maxLines {
		if budget := l.rowsAvailable(twoColumnChrome); budget > 0 && detail > budget {
			detail = budget
		}
		if detail > maxLines {
			maxLines = detail
		}
	}

	for i := 0; i < maxLines; i++ {
		leftLine := ""
		if i < len(listLines) {
			leftLine = listLines[i]
		}
		rightLine := ""
		if i < len(detailLines) {
			rightLine = detailLines[i]
		}

		leftLine = lipgloss.NewStyle().Width(leftWidth).Render(leftLine)

		b.WriteString(leftLine)
		b.WriteString(" │ ")
		b.WriteString(rightLine)
		b.WriteString("\n")
	}

	b.WriteString("\n")

	if l.Message != "" {
		msgStyle := lipgloss.NewStyle().Foreground(gold).Bold(true)
		b.WriteString(msgStyle.Render(i18n.T("success_prefix") + " " + l.Message))
		b.WriteString("\n")
	}

	help := components.NewHelpBar()
	help.Shortcuts = l.Shortcuts
	help.Width = l.Width
	b.WriteString(help.Render())

	return b.String()
}

func (l *ListView) singleColumn() string {
	var b strings.Builder

	gold := ui.ThemeDefault.Gold
	textDim := ui.ThemeDefault.TextDim

	title := lipgloss.NewStyle().
		Foreground(gold).
		Bold(true).
		Render(i18n.T("app_name_tui"))
	versionStr := lipgloss.NewStyle().
		Foreground(textDim).
		Render("v" + version.Version + "  ws:" + fmt.Sprintf("%d", l.Sessions))

	b.WriteString(title)
	b.WriteString(ui.Repeat(" ", l.Width-lipgloss.Width(title)-lipgloss.Width(versionStr)-ui.HeaderMargin))
	b.WriteString(versionStr)
	b.WriteString("\n")

	if l.UpdateBadge != "" {
		badgeStyle := lipgloss.NewStyle().Foreground(ui.ThemeDefault.Gold).Bold(true)
		b.WriteString(badgeStyle.Render(l.UpdateBadge))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	l.ListTop = 2
	if l.UpdateBadge != "" {
		l.ListTop++
	}
	if l.Dashboard != "" {
		l.ListTop += 2
	}
	l.VisibleRows = l.rowsAvailable(singleColumnChrome)

	if l.Dashboard != "" {
		urlStyle := lipgloss.NewStyle().Foreground(gold)
		b.WriteString(urlStyle.Render(i18n.T("dashboard_label") + " " + l.Dashboard + " " + i18n.T("press_w_hint")))
		b.WriteString("\n\n")
	}

	b.WriteString(l.renderTunnelList(l.Width - 2))
	b.WriteString("\n")

	if l.Message != "" {
		msgStyle := lipgloss.NewStyle().Foreground(gold).Bold(true)
		b.WriteString(msgStyle.Render(i18n.T("success_prefix") + " " + l.Message))
		b.WriteString("\n")
	}

	help := components.NewHelpBar()
	help.Shortcuts = l.Shortcuts
	help.Width = l.Width
	b.WriteString(help.Render())

	return b.String()
}

func (l *ListView) visibleWindow() (first, count int) {
	total := len(l.Items)
	if l.VisibleRows <= 0 || l.VisibleRows >= total {
		return 0, total
	}

	first = l.Cursor - l.VisibleRows/2
	if first < 0 {
		first = 0
	}
	if first > total-l.VisibleRows {
		first = total - l.VisibleRows
	}

	return first, l.VisibleRows
}

func (l *ListView) renderTunnelList(width int) string {
	var b strings.Builder

	first, count := l.visibleWindow()
	l.FirstItem = first
	last := first + count

	if first > 0 {
		b.WriteString(l.scrollHint(width, i18n.TF("more_above", first)))
		b.WriteString("\n")
		l.ListTop++
	}

	for i := first; i < last; i++ {
		item := l.Items[i]
		tunnelItem := components.TunnelItem{
			Selected:    i == l.Cursor,
			Name:        item.Name,
			Provider:    item.Provider,
			LocalPort:   item.LocalPort,
			StatusState: item.StatusState,
			StatusMsg:   item.StatusMsg,
			Width:       width,
		}
		b.WriteString(tunnelItem.Render())
		if i < last-1 {
			b.WriteString("\n")
		}
	}

	if remaining := len(l.Items) - last; remaining > 0 {
		b.WriteString("\n")
		b.WriteString(l.scrollHint(width, i18n.TF("more_below", remaining)))
	}

	return b.String()
}

func (l *ListView) scrollHint(width int, text string) string {
	return lipgloss.NewStyle().
		Foreground(ui.ThemeDefault.TextDim).
		Width(ui.Clamp(width, 0)).
		Align(lipgloss.Center).
		Render(text)
}

func (l *ListView) renderDetailPanel(width int) string {
	if l.Cursor < 0 || l.Cursor >= len(l.Items) {
		return lipgloss.NewStyle().
			Foreground(ui.ThemeDefault.TextDim).
			Render(i18n.T("select_tunnel_details"))
	}

	item := l.Items[l.Cursor]

	panel := components.DetailPanel{
		Name:        item.Name,
		Provider:    item.Provider,
		LocalPort:   item.LocalPort,
		StatusState: item.StatusState,
		StatusMsg:   item.StatusMsg,
		PublicURL:   item.PublicURL,
		ErrorMsg:    item.ErrorMsg,
		Width:       width - 2,
	}

	return panel.Render()
}
