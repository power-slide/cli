package updater

import (
	"fmt"
	"net/http"
	"time"

	"github.com/inconshreveable/go-update"
	"github.com/power-slide/cli/pkg/config"

	"github.com/power-slide/cli/pkg/logger"
	"github.com/power-slide/cli/pkg/version"
	log "github.com/sirupsen/logrus"
	"golang.org/x/mod/semver"
)

func CheckNeeded() bool {
	nextCheck := config.LastUpdateCheck().Add(config.AutoUpdateInterval())
	return config.AutoUpdateEnabled() && version.IsReleaseBuild() && nextCheck.Before(time.Now())
}

func UpdateNeeded() bool {
	if CheckNeeded() {
		config.SetLastUpdateCheck()
		config.WriteConfig()

		localVersion := fmt.Sprintf("v%s", version.Version)
		if !semver.IsValid(localVersion) {
			return true
		} else {
			remoteVersion := fmt.Sprintf("v%s", version.LatestVersion())
			if semver.Compare(localVersion, remoteVersion) == -1 {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func AutomaticUpdate() {
	if !UpdateNeeded() {
		log.Debugln("Skipped auto update")
		return
	}

	var answer string = "noop"
	answers := []string{"y", "n", ""}
	var doBreak bool
	for {
		for _, a := range answers {
			if a == answer {
				doBreak = true
			}
		}
		if doBreak {
			break
		}
		fmt.Print("Update available, would you like to skip running your command to update? (Y/n) ")
		fmt.Scanln(&answer)
	}

	if answer == "y" || answer == "" {
		Update()
		fmt.Println("Exiting, please re-run your command.")
		log.Exit(0)
	}
}

func Update() {
	log.Debugln("Updating CLI binary")

	url := version.LatestVersionURL()
	if url == "" {
		logger.CheckErr(fmt.Errorf("unable to get release URL"))
	}

	var err error
	resp, err := http.Get(url)
	logger.CheckErr(err)
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	logger.CheckErr(err)
	log.Debugln("Finished the update")
}
