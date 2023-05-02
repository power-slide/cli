package updater

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/google/go-github/github"
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

	release, err := version.LatestRelease()
	logger.CheckErr(err)

	url := assetURL(release, getTargetName())
	if url == "" {
		logger.CheckErr(fmt.Errorf("unable to get release URL"))
	}

	resp, err := http.Get(url)
	logger.CheckErr(err)
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{Checksum: getChecksum(release)})
	logger.CheckErr(err)
	log.Debugln("Finished the update")
}

func getTargetName() string {
	targetName := fmt.Sprintf("pwrsl-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		targetName = fmt.Sprintf("%s.exe", targetName)
	}
	return targetName
}

func assetURL(release *github.RepositoryRelease, asset string) string {
	var url string
	for _, a := range release.Assets {
		if *a.Name == asset {
			url = *a.BrowserDownloadURL
			break
		}
	}
	return url
}

func getChecksum(release *github.RepositoryRelease) []byte {
	url := assetURL(release, "sha256sums")
	if url == "" {
		logger.CheckErr(fmt.Errorf("unable to get checksum URL"))
	}

	resp, err := http.Get(url)
	logger.CheckErr(err)
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.CheckErr(err)
	}

	checksums := strings.Split(string(responseData), "\n")
	targetName := getTargetName()
	var hexChecksum string
	for _, archCheckSum := range checksums {
		if strings.Contains(archCheckSum, targetName) {
			hexChecksum = strings.Split(archCheckSum, " ")[0]
			break
		}
	}

	if hexChecksum == "" {
		logger.CheckErr(fmt.Errorf("unable to get checksum"))
	}

	checksum, err := hex.DecodeString(hexChecksum)
	logger.CheckErr(err)
	return checksum
}
