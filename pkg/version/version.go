package version

import (
	"strings"
)

var (
	Version string
)

func IsReleaseBuild() bool {
	return !strings.HasPrefix(Version, "dev.")
}
