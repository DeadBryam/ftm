package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"

	"github.com/sthbryan/ftm/internal/app/ui/views"
	"github.com/sthbryan/ftm/internal/config"
	"github.com/sthbryan/ftm/internal/i18n"
	"github.com/sthbryan/ftm/internal/providers"
	"github.com/sthbryan/ftm/internal/updater"
)

type viewState int

const (
	viewList viewState = iota
	viewLogs
	viewAddForm
	viewEditForm
	viewConfirm
	viewDownloading
	viewSettings
)

type Settings struct {
	NotificationsEnabled bool
	NotificationSound    bool
	Theme                string
}

const TwoColumnThreshold = 100

type Model struct {
	App                 *App
	Keys                KeyMap
	Help                help.Model
	State               viewState
	Width               int
	Height              int
	Cursor              int
	Items               []list.Item
	LogViewport         viewport.Model
	SelectedTunnel      string
	FormFocus           int
	FormValues          FormData
	editingTunnelID     string
	Message             string
	MessageTimer        int
	DownloadProgress    providers.DownloadProgress
	DownloadingProvider string
	PendingTunnel       *config.TunnelConfig
	ProgressBar         progress.Model
	SettingsView        *views.SettingsView
	UpdateAvailable     *updater.Info
}

type FormData struct {
	ID       string
	Name     string
	Provider string
	Port     string
}

type TunnelItem struct {
	Tunnel config.TunnelConfig
	Status config.TunnelStatus
}

func (i TunnelItem) FilterValue() string { return i.Tunnel.Name }

func (i TunnelItem) Title() string { return i.Tunnel.Name }

func (i TunnelItem) Description() string {
	status := i18n.T("status_offline")
	switch i.Status.State {
	case config.TunnelStateStarting:
		status = i18n.T("status_starting")
	case config.TunnelStateConnecting:
		status = i18n.T("status_connecting")
	case config.TunnelStateOnline:
		status = i18n.T("status_online")
	case config.TunnelStateError:
		status = i18n.T("status_error")
	case config.TunnelStateTimeout:
		status = i18n.T("status_timeout")
	}
	return fmt.Sprintf("%s | %s %d | %s", i.Tunnel.Provider, i18n.T("port"), i.Tunnel.LocalPort, status)
}
