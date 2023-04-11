package logger

import (
	"os"
	"strings"

	"github.com/power-slide/cli/pkg/version"
	log "github.com/sirupsen/logrus"
)

func Init() {
	log.SetOutput(os.Stderr)

	if !strings.HasPrefix(version.Version, "test") && version.IsReleaseBuild() {
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
}

func CheckErr(msg interface{}) {
	if msg != nil {
		log.Errorln("Error:", msg)
		log.Exit(1)
	}
}
