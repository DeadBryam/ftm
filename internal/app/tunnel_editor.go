package app

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/sthbryan/ftm/internal/config"
	"github.com/sthbryan/ftm/internal/i18n"
)

func (m *Model) handleEditorKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.EditorFocus = (m.EditorFocus + 1) % 5

	case "shift+tab":
		m.EditorFocus = (m.EditorFocus - 1 + 5) % 5

	case "enter":
		return m.handleEditorEnter()

	case "esc":
		m.State = viewList
		m.editingTunnelID = ""

	case "left", "right":
		m.handleProviderNav(msg.String())

	default:
		m.handleEditorInput(msg.String())
	}

	return m, nil
}

func (m *Model) handleEditorEnter() (tea.Model, tea.Cmd) {
	switch m.EditorFocus {
	case 4:
		m.saveTunnel()
	default:
		m.EditorFocus++
	}
	return m, nil
}

func (m *Model) handleProviderNav(dir string) {
	if m.EditorFocus != 1 {
		return
	}

	providers := []config.Provider{
		config.ProviderCloudflared,
		config.ProviderTunnelmole,
		config.ProviderLocalhostRun,
		config.ProviderServeo,
		config.ProviderPinggy,
		config.ProviderBore,
	}

	current := config.Provider(m.Draft.Provider)
	idx := -1
	for i, p := range providers {
		if p == current {
			idx = i
			break
		}
	}

	if idx == -1 {
		m.Draft.Provider = string(config.ProviderCloudflared)
		return
	}

	if dir == "right" {
		idx = (idx + 1) % len(providers)
	} else {
		idx = (idx - 1 + len(providers)) % len(providers)
	}

	m.Draft.Provider = string(providers[idx])
}

func (m *Model) handleEditorInput(s string) {
	switch m.EditorFocus {
	case 0:
		m.handleNameInput(s)

	case 2:
		m.handlePortInput(s)
	}
}

func (m *Model) handleNameInput(s string) {
	if s == "backspace" {
		if len(m.Draft.Name) > 0 {
			m.Draft.Name = m.Draft.Name[:len(m.Draft.Name)-1]
		}
	} else if s == "space" {
		m.Draft.Name += " "
	} else if len(s) == 1 {
		m.Draft.Name += s
	}
}

func (m *Model) handlePortInput(s string) {
	if s == "backspace" {
		if len(m.Draft.Port) > 0 {
			m.Draft.Port = m.Draft.Port[:len(m.Draft.Port)-1]
		}
	} else if s >= "0" && s <= "9" && len(m.Draft.Port) < 5 {
		m.Draft.Port += s
	}
}

func (m *Model) saveTunnel() {
	if m.Draft.Name == "" || m.Draft.Port == "" {
		m.showMessage(i18n.T("validation_required_fields"))
		return
	}

	if m.editingTunnelID != "" {
		m.updateTunnel()
	} else {
		m.createTunnel()
	}
}

func (m *Model) updateTunnel() {
	for i := range m.App.Config.Tunnels {
		if m.App.Config.Tunnels[i].ID == m.editingTunnelID {
			m.App.Config.Tunnels[i].Name = m.Draft.Name
			m.App.Config.Tunnels[i].Provider = config.Provider(m.Draft.Provider)
			m.App.Config.Tunnels[i].LocalPort = parsePort(m.Draft.Port)

			if m.App.WebServer != nil {
				m.App.WebServer.BroadcastTunnelUpdate(m.App.Config.Tunnels[i])
			}
			break
		}
	}
	m.editingTunnelID = ""
	m.App.SaveConfig()
	m.refreshItems()
	m.State = viewList
	m.showMessage(i18n.T("tunnel_updated"))
}

func (m *Model) createTunnel() {
	id := strings.ToLower(strings.ReplaceAll(m.Draft.Name, " ", "-"))

	tunnel := config.TunnelConfig{
		ID:        id,
		Name:      m.Draft.Name,
		Provider:  config.Provider(m.Draft.Provider),
		LocalPort: parsePort(m.Draft.Port),
	}

	m.App.Config.AddTunnel(tunnel)
	m.App.SaveConfig()
	m.refreshItems()
	m.State = viewList
	m.showMessage(i18n.T("tunnel_added"))
}

func parsePort(s string) int {
	var port int
	for _, c := range s {
		port = port*10 + int(c-'0')
	}
	return port
}
