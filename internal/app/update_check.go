package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/sthbryan/ftm/internal/updater"
	"github.com/sthbryan/ftm/internal/version"
)

const updateRepo = "sthbryan/ftm"

type updateCheckMsg struct {
	info *updater.Info
	err  error
}

type updateApplyMsg struct {
	err error
}

func checkUpdateCmd() tea.Cmd {
	return func() tea.Msg {
		info, err := updater.New(updateRepo).Check(version.Version)
		return updateCheckMsg{info: info, err: err}
	}
}

func applyUpdateCmd(info *updater.Info) tea.Cmd {
	return func() tea.Msg {
		err := updater.New(updateRepo).Apply(info)
		return updateApplyMsg{err: err}
	}
}
