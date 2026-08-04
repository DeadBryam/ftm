package cli

import (
	"github.com/sthbryan/ftm/internal/config"
	"github.com/sthbryan/ftm/internal/i18n"
)

//

func Init() error {
	if err := i18n.Load(); err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {

		i18n.SetLanguageWithFallback(i18n.ResolveLanguage(""))
		return nil
	}

	i18n.InitFromConfig(cfg)

	return nil
}
