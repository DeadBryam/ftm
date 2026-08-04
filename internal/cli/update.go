package cli

import (
	"fmt"

	"github.com/sthbryan/ftm/internal/i18n"
	"github.com/sthbryan/ftm/internal/updater"
	"github.com/sthbryan/ftm/internal/version"
)

const defaultRepo = "sthbryan/ftm"

func Update(checkOnly bool) error {
	u := updater.New(defaultRepo)
	info, err := u.Check(version.Version)
	if err != nil {
		return fmt.Errorf("%s", i18n.TF("update_check_failed", err.Error()))
	}

	fmt.Println(i18n.TF("update_current_version", info.CurrentVersion))
	fmt.Println(i18n.TF("update_latest_version", info.LatestVersion))

	if !info.HasUpdate {
		fmt.Println(i18n.T("update_up_to_date"))
		return nil
	}

	fmt.Println(i18n.TF("update_available", info.Tag))
	fmt.Println(i18n.TF("update_release_url", info.ReleaseURL))
	if checkOnly {
		return nil
	}

	if info.Method != updater.MethodSelf {
		fmt.Println(i18n.T(info.Method.HintKey()))
		return nil
	}

	fmt.Println(i18n.TF("update_downloading", info.AssetName))
	if err := u.Apply(info); err != nil {
		return fmt.Errorf("%s", i18n.TF("update_apply_failed", err.Error()))
	}
	fmt.Println(i18n.TF("update_success", info.LatestVersion))
	return nil
}
