package views

import (
	"strings"
	"testing"
)

func manyItems(n int) []TunnelViewData {
	items := make([]TunnelViewData, 0, n)
	for i := 0; i < n; i++ {
		items = append(items, TunnelViewData{
			Name:        "Tunnel",
			Provider:    "cloudflared",
			LocalPort:   30000 + i,
			StatusState: 3,
		})
	}
	return items
}

func listWith(items []TunnelViewData, width, height, cursor int) *ListView {
	v := NewListView()
	v.Width = width
	v.Height = height
	v.Items = items
	v.Cursor = cursor
	v.Dashboard = "http://localhost:40500"
	v.Shortcuts = nil
	return v
}

func TestListNeverRendersTallerThanTheTerminal(t *testing.T) {
	for _, height := range []int{10, 16, 24, 40} {
		for _, width := range []int{60, 120} {
			v := listWith(manyItems(50), width, height, 0)
			out := v.Render()

			if lines := strings.Count(out, "\n") + 1; lines > height {
				t.Errorf("width %d height %d: rendered %d lines", width, height, lines)
			}
		}
	}
}

func TestSelectedItemStaysVisibleWhileScrolling(t *testing.T) {
	items := manyItems(50)

	for _, cursor := range []int{0, 7, 25, 49} {
		v := listWith(items, 60, 20, cursor)
		v.Render()

		first, count := v.visibleWindow()
		if cursor < first || cursor >= first+count {
			t.Errorf("cursor %d falls outside the visible window [%d,%d)", cursor, first, first+count)
		}
	}
}

func TestShortListIsNotWindowed(t *testing.T) {
	v := listWith(manyItems(3), 60, 40, 1)
	v.Render()

	first, count := v.visibleWindow()
	if first != 0 || count != 3 {
		t.Errorf("window = [%d,%d), want the whole list", first, first+count)
	}
	if strings.Contains(v.Render(), "more") {
		t.Error("a list that fits still showed a scroll hint")
	}
}

func TestScrollHintsAppearOnlyWhenThereIsMore(t *testing.T) {
	items := manyItems(50)

	top := listWith(items, 60, 20, 0)
	outTop := top.Render()
	if strings.Contains(outTop, "↑") {
		t.Error("an 'above' hint showed while at the top of the list")
	}
	if !strings.Contains(outTop, "↓") {
		t.Error("no 'below' hint while items remain below")
	}

	bottom := listWith(items, 60, 20, len(items)-1)
	outBottom := bottom.Render()
	if !strings.Contains(outBottom, "↑") {
		t.Error("no 'above' hint while items remain above")
	}
	if strings.Contains(outBottom, "↓") {
		t.Error("a 'below' hint showed while at the bottom of the list")
	}
}

func TestVisibleWindowKeepsItsSize(t *testing.T) {
	items := manyItems(50)

	for _, cursor := range []int{0, 10, 25, 49} {
		v := listWith(items, 60, 20, cursor)
		v.Render()

		first, count := v.visibleWindow()
		if first < 0 {
			t.Fatalf("cursor %d produced a negative window start", cursor)
		}
		if first+count > len(items) {
			t.Fatalf("cursor %d produced a window past the end: [%d,%d)", cursor, first, first+count)
		}
	}
}

func TestTwoColumnKeepsTheScrollHint(t *testing.T) {
	v := listWith(manyItems(50), 120, 16, 0)
	out := v.Render()

	if !strings.Contains(out, "↓") {
		t.Error("the two-column layout dropped the 'more below' hint")
	}
}

func TestTwoColumnDetailPanelDoesNotOverflow(t *testing.T) {
	items := manyItems(2)
	items[0].PublicURL = "https://random-words-here.trycloudflare.com"
	items[0].ErrorMsg = "a fairly long error message that wraps"

	for _, height := range []int{12, 16, 24} {
		v := listWith(items, 120, height, 0)
		out := v.Render()

		if lines := strings.Count(out, "\n") + 1; lines > height {
			t.Errorf("height %d: rendered %d lines", height, lines)
		}
	}
}

func TestLogsViewNeverRendersTallerThanTheTerminal(t *testing.T) {
	var content strings.Builder
	for i := 0; i < 500; i++ {
		content.WriteString("2026-08-03 12:00:00 starting tunnel\n")
	}

	for _, height := range []int{8, 16, 24, 40} {
		v := NewLogsView()
		v.Width, v.Height = 100, height
		v.TunnelName = "Foundry VTT"
		v.Content = content.String()

		if lines := strings.Count(v.Render(), "\n") + 1; lines > height {
			t.Errorf("height %d: rendered %d lines", height, lines)
		}
	}
}

func TestLogsViewShowsTheNewestLines(t *testing.T) {
	v := NewLogsView()
	v.Width, v.Height = 100, 12
	v.Content = "first line\nmiddle line\nnewest line\n"

	out := v.Render()
	if !strings.Contains(out, "newest line") {
		t.Error("the newest log line is not visible")
	}
}

func TestCentredViewsFillTheTerminal(t *testing.T) {
	const width, height = 100, 30

	form := NewFormView()
	form.Width, form.Height = width, height

	empty := NewEmptyState()
	empty.Width, empty.Height = width, height

	settings := NewSettingsView()
	settings.Width, settings.Height = width, height

	for name, out := range map[string]string{
		"form": form.Render(), "empty state": empty.Render(), "settings": settings.Render(),
	} {
		if lines := strings.Count(out, "\n") + 1; lines != height {
			t.Errorf("%s: rendered %d lines, want %d", name, lines, height)
		}
	}
}
