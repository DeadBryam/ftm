package app

import (
	"github.com/charmbracelet/bubbles/key"

	"github.com/sthbryan/ftm/internal/app/ui/components"
	"github.com/sthbryan/ftm/internal/i18n"
)

// KeyMap is the single source of truth for keyboard shortcuts.
//
// The help bar and the README used to hardcode their own lists, and all three
// disagreed: the README documented "s" for start/stop and "l" for logs, while
// the code bound "s" to settings, and "l" to both logs and an unused "next"
// action.
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	Toggle   key.Binding
	Logs     key.Binding
	Copy     key.Binding
	Web      key.Binding
	Add      key.Binding
	Edit     key.Binding
	Delete   key.Binding
	Config   key.Binding
	Settings key.Binding
	Update   key.Binding
	Back     key.Binding
	Help     key.Binding
	Quit     key.Binding

	// ForceQuit is the interrupt, and is the only shortcut that still applies
	// while a text field has focus.
	ForceQuit key.Binding
}

var DefaultKeys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	// Arrows only: "h"/"l" would collide with the logs shortcut, which is what
	// the previous map did.
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "previous"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "next"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "start/stop"),
	),
	Logs: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "logs"),
	),
	Copy: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy URL"),
	),
	Web: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "open web"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Config: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open config"),
	),
	Settings: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "settings"),
	),
	Update: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "update"),
	),
	// Esc belongs to Back alone. It used to be bound to Quit as well, and
	// because Quit is matched first, Esc closed the app from the list view.
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	ForceQuit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Toggle},
		{k.Logs, k.Copy, k.Web},
		{k.Add, k.Edit, k.Delete},
		{k.Settings, k.Config, k.Update},
		{k.Back, k.Help, k.Quit},
	}
}

// listShortcuts is what the help bar under the tunnel list shows. The labels
// are translated; the keys come from the bindings above so the two cannot drift
// apart.
func (k KeyMap) listShortcuts() []components.Shortcut {
	return []components.Shortcut{
		// Navigation spans two bindings, so it is the one entry whose keys are
		// written out rather than taken from a single binding.
		{Keys: "↑↓/kj", Label: i18n.T("navigate")},
		{Keys: k.Enter.Help().Key + "/" + k.Toggle.Help().Key, Label: i18n.T("start_stop")},
		{Keys: k.Logs.Help().Key, Label: i18n.T("logs")},
		{Keys: k.Copy.Help().Key, Label: i18n.T("copy_url")},
		{Keys: k.Add.Help().Key, Label: i18n.T("create")},
		{Keys: k.Edit.Help().Key, Label: i18n.T("edit")},
		{Keys: k.Delete.Help().Key, Label: i18n.T("delete")},
		{Keys: k.Settings.Help().Key, Label: i18n.T("settings")},
		{Keys: k.Web.Help().Key, Label: i18n.T("web")},
		{Keys: k.Config.Help().Key, Label: i18n.T("config")},
		{Keys: k.Quit.Help().Key, Label: i18n.T("close")},
	}
}
