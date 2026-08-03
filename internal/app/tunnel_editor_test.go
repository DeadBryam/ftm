package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func editorWith(t *testing.T, focus int) *Model {
	t.Helper()

	m := modelWithTunnel(t, "t1", "Foundry VTT")
	m.State = viewNewTunnel
	m.EditorFocus = focus

	return m
}

func TestArrowsMoveBetweenFields(t *testing.T) {
	m := editorWith(t, 0)

	m.handleKey(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.EditorFocus != 1 {
		t.Errorf("down left the focus on %d, want 1", m.EditorFocus)
	}

	m.handleKey(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.EditorFocus != 0 {
		t.Errorf("up left the focus on %d, want 0", m.EditorFocus)
	}
}

func TestFocusWrapsAroundAndNeverStopsOnNothing(t *testing.T) {
	m := editorWith(t, 0)

	seen := make(map[int]bool)
	for i := 0; i < editorStops; i++ {
		seen[m.EditorFocus] = true
		m.handleKey(tea.KeyPressMsg{Code: tea.KeyDown})
	}

	if m.EditorFocus != 0 {
		t.Errorf("a full cycle ended on %d, want back at 0", m.EditorFocus)
	}
	for stop := 0; stop < editorStops; stop++ {
		if !seen[stop] {
			t.Errorf("focus %d is never reachable", stop)
		}
	}
}

func TestUpFromTheFirstFieldWrapsToSubmit(t *testing.T) {
	m := editorWith(t, 0)

	m.handleKey(tea.KeyPressMsg{Code: tea.KeyUp})

	if m.EditorFocus != editorSubmit {
		t.Errorf("focus = %d, want the submit button at %d", m.EditorFocus, editorSubmit)
	}
}

func TestArrowsStillChangeTheProvider(t *testing.T) {
	m := editorWith(t, 1)
	m.Draft.Provider = "cloudflared"

	m.handleKey(tea.KeyPressMsg{Code: tea.KeyRight})

	if m.Draft.Provider == "cloudflared" {
		t.Error("right did not move to the next provider")
	}
	if m.EditorFocus != 1 {
		t.Errorf("changing the provider moved the focus to %d", m.EditorFocus)
	}
}
