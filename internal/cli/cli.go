package cli

import (
	"github.com/sthbryan/ftm/internal/config"
	"github.com/sthbryan/ftm/internal/i18n"
)

// Init loads translations and selects the UI language before any command runs.
//
// Language selection belongs here rather than only in app.New: --version,
// --update and --uninstall all exit before an App is ever built, so they used
// to print English no matter what the user had configured.
func Init() error {
	if err := i18n.Load(); err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		// An unreadable config must not stop --version or --uninstall from
		// working. Fall back to the system language and carry on; the commands
		// that actually need a config will report the error themselves.
		i18n.SetLanguageWithFallback(i18n.ResolveLanguage(""))
		return nil
	}

	i18n.InitFromConfig(cfg)

	return nil
}
