package update

import (
	"fmt"
	"time"

	"github.com/power-slide/cli/pkg/config"
	"github.com/power-slide/cli/pkg/updater"
	"github.com/power-slide/cli/pkg/version"
	"github.com/spf13/cobra"
	str2dur "github.com/xhit/go-str2duration/v2"
)

var (
	Cmd = &cobra.Command{
		Use:   "update",
		Short: "Update this CLI to the latest version",
		Long:  "Update this CLI to the latest version and manage the automatic update settings.",
		Run:   Run,
		Args:  checkArgs,
	}

	skipMainAction   bool
	forceUpdate      bool
	toggleAutoUpdate bool
	setAutoUpdate    string
	setCheckInterval string
	checkInterval    time.Duration
)

func init() {
	Cmd.Flags().BoolVar(&toggleAutoUpdate, "toggle-auto-update", false, "Enable/Disable automatic updating of this CLI")
	Cmd.Flags().StringVar(&setAutoUpdate, "auto-update", "", "Set to true/false to enable/disable")
	Cmd.Flags().StringVar(&setCheckInterval, "check-interval", "", "Set to a duration between 1h and 7d")
	Cmd.Flags().BoolVarP(&forceUpdate, "force", "f", false, "Force an update, disregard checking interval")
	Cmd.MarkFlagsMutuallyExclusive("toggle-auto-update", "auto-update")
}

func checkArgs(cmd *cobra.Command, args []string) error {
	var validUpdatevalue bool
	for _, v := range []string{"", "true", "false"} {
		if v == setAutoUpdate {
			validUpdatevalue = true
		}
	}

	if !validUpdatevalue {
		return fmt.Errorf("auto-update must be either true or false")
	}

	if setCheckInterval != "" {
		var err error
		checkInterval, err = str2dur.ParseDuration(setCheckInterval)
		if err != nil {
			return fmt.Errorf("invalid check-interval value: %w", err)
		}
	}

	return nil
}

func Run(cmd *cobra.Command, _ []string) {
	handleAutoUpdateToggle()
	handleSetAutoUpdate()
	handleSetCheckInterval()

	performUpdate()
}

func handleAutoUpdateToggle() {
	if !toggleAutoUpdate {
		return
	}

	skipMainAction = true
	config.ToggleAutoUpdate()
	if config.AutoUpdateEnabled() {
		fmt.Println("Auto-update enabled")
	} else {
		fmt.Println("Auto-update disabled")
	}
}

func handleSetAutoUpdate() {
	if setAutoUpdate == "" {
		return
	} else if setAutoUpdate == "true" {
		config.SetAutoUpdate(true)
		fmt.Println("Auto-update enabled")
	} else {
		config.SetAutoUpdate(false)
		fmt.Println("Auto-update disabled")
	}
	skipMainAction = true
}

func handleSetCheckInterval() {
	if setCheckInterval == "" {
		return
	}

	min, _ := str2dur.ParseDuration("1h")
	max, _ := str2dur.ParseDuration("1w")

	if checkInterval < min || checkInterval > max {
		fmt.Println("Update check interval", str2dur.String(checkInterval), "must be between 1h and 1w, config not updated")
	} else {
		config.SetAutoUpdateInterval(checkInterval)
		fmt.Println("Update check interval set to", str2dur.String(checkInterval))
	}

	skipMainAction = true
}

func performUpdate() {
	if skipMainAction {
		return
	}

	if forceUpdate || updater.UpdateNeeded() {
		fmt.Print("Performing update... ")
		updater.Update()
		fmt.Println("Done!")
	} else {
		fmt.Println("Up to date!")
		fmt.Println("Latest Version:", version.LatestVersion())
		fmt.Println("Local Version:", version.Version)
	}
}
