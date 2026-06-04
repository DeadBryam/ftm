package cli

import "github.com/sthbryan/ftm/internal/i18n"

func Init() error {
	return i18n.Load()
}
