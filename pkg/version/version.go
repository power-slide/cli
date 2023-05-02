package version

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

var (
	Version string
)

func IsReleaseBuild() bool {
	return !strings.HasPrefix(Version, "dev.")
}

func LatestRelease() (*github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	release, resp, err := client.Repositories.GetLatestRelease(ctx, "power-slide", "cli")
	if resp.StatusCode >= 400 {
		return release, fmt.Errorf("release not found")
	}

	if ctx.Err() != nil {
		err = ctx.Err()
	}

	return release, err
}

func LatestVersion() string {
	release, err := LatestRelease()
	if err != nil {
		return "unknown"
	} else {
		return strings.ReplaceAll(*release.TagName, "v", "")
	}
}

func LatestVersionURL() string {
	targetName := fmt.Sprintf("pwrsl-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		targetName = fmt.Sprintf("%s.exe", targetName)
	}
	var url string
	release, _ := LatestRelease()

	for _, a := range release.Assets {
		if *a.Name == targetName {
			log.Debugln("Latest release download url:", *a.BrowserDownloadURL)
			url = *a.BrowserDownloadURL
			break
		}
	}

	return url
}
