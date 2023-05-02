package version

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
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
