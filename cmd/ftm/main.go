package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sthbryan/ftm/internal/app"
	"github.com/sthbryan/ftm/internal/cli"
	"github.com/sthbryan/ftm/internal/i18n"
	"github.com/sthbryan/ftm/internal/updater"
	"github.com/sthbryan/ftm/internal/version"
)

var BuildVersion string

func doUpdate(checkOnly bool) {
	_ = i18n.Load()
	u := updater.New("sthbryan/ftm")
	info, err := u.Check(version.Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, i18n.TF("update_check_failed", err.Error())+"\n")
		os.Exit(1)
	}

	fmt.Println(i18n.TF("update_current_version", info.CurrentVersion))
	fmt.Println(i18n.TF("update_latest_version", info.LatestVersion))

	if !info.HasUpdate {
		fmt.Println(i18n.T("update_up_to_date"))
		return
	}

	fmt.Println(i18n.TF("update_available", info.Tag))
	fmt.Println(i18n.TF("update_release_url", info.ReleaseURL))
	if checkOnly {
		return
	}

	fmt.Println(i18n.TF("update_downloading", info.AssetName))
	if err := u.Apply(info); err != nil {
		fmt.Fprintf(os.Stderr, i18n.TF("update_apply_failed", err.Error())+"\n")
		os.Exit(1)
	}
	fmt.Println(i18n.TF("update_success", info.LatestVersion))
}

func main() {
	_ = cli.Init()
	var (
		webOnly     = flag.Bool("web", false, "Start web dashboard and open browser")
		server      = flag.Bool("server", false, "Start web dashboard only (no browser)")
		port        = flag.Int("port", 0, "Web server port (auto-detect if not specified)")
		showVersion = flag.Bool("version", false, "Show version")
		uninstall   = flag.Bool("uninstall", false, "Uninstall ftm")
		update      = flag.Bool("update", false, "Update ftm to the latest release")
		checkOnly   = flag.Bool("check", false, "Check for updates without installing")
	)
	flag.Parse()

	if *showVersion {
		fmt.Println(i18n.TF("version_output", version.Version))
		os.Exit(0)
	}

	if *uninstall {
		if err := cli.Uninstall(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *checkOnly {
		doUpdate(true)
		os.Exit(0)
	}

	if *update {
		doUpdate(false)
		os.Exit(0)
	}

	if BuildVersion == "" {
		BuildVersion = version.Version
	}

	application, err := app.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *port > 0 {
		application.Config.WebPort = *port
	}

	if err := application.StartWebServer(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting web server: %v\n", err)
		os.Exit(1)
	}

	url := application.WebServer.URL()
	fmt.Printf("🎲 Foundry Tunnel Manager v%s\n", BuildVersion)
	fmt.Printf(i18n.TF("dashboard_url", url))

	if *webOnly {
		fmt.Print(i18n.T("press_ctrl_c"))
		application.OpenDashboard()
		select {}
	} else if *server {
		fmt.Print(i18n.T("press_ctrl_c"))
		select {}
	}

	fmt.Print(i18n.T("tui_hint"))

	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
