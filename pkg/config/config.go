package config

import (
	"os"
	"time"

	"github.com/power-slide/cli/pkg/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	str2dur "github.com/xhit/go-str2duration/v2"
)

const (
	DefaultConfigFile   = ".powerslide"
	DefaultConfigFormat = "yaml"
	autoUpdateKey       = "auto-update"
	updateIntervalKey   = "update-interval"
	lastUpdateCheckKey  = "last-update-check"

	defaultLastUpdate     = 0
	defaultAutoUpdate     = true
	defaultUpdateInterval = "1d"
)

var (
	ConfigFile     string
	configUpdated  bool
	existingConfig bool
)

func Init() {
	if ConfigFile != "" {
		viper.SetConfigFile(ConfigFile)
	} else {
		home, err := os.UserHomeDir()
		logger.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType(DefaultConfigFormat)
		viper.SetConfigName(DefaultConfigFile)
	}

	setDefaultConfig()
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			existingConfig = false
			configUpdated = true
		} else {
			logger.CheckErr(err)
		}
	} else {
		log.WithFields(rawFields()).Debugln("config: loaded from", viper.ConfigFileUsed())
		existingConfig = true
		configUpdated = false
	}
}

func setDefaultConfig() {
	viper.SetDefault(autoUpdateKey, defaultAutoUpdate)
	viper.SetDefault(updateIntervalKey, defaultUpdateInterval)
	viper.SetDefault(lastUpdateCheckKey, defaultLastUpdate)
}

func AutoUpdateEnabled() bool {
	return viper.GetBool(autoUpdateKey)
}

func SetAutoUpdate(newValue bool) {
	configUpdated = true
	log.Debugln("config: setting", autoUpdateKey, "to", newValue)
	viper.Set(autoUpdateKey, newValue)
}

func ToggleAutoUpdate() {
	newValue := !viper.GetBool(autoUpdateKey)
	configUpdated = true
	log.Debugln("config: toggling", autoUpdateKey, "to", newValue)
	viper.Set(autoUpdateKey, newValue)
}

func AutoUpdateInterval() time.Duration {
	interval, err := str2dur.ParseDuration(viper.GetString(updateIntervalKey))
	if err != nil {
		interval, _ = str2dur.ParseDuration(defaultUpdateInterval)
	}
	return interval
}

func SetAutoUpdateInterval(interval time.Duration) {
	newValue := str2dur.String(interval)
	configUpdated = true
	log.Debugln("config: setting", updateIntervalKey, "to", newValue)
	viper.Set(updateIntervalKey, newValue)
}

func LastUpdateCheck() time.Time {
	return time.Unix(viper.GetInt64(lastUpdateCheckKey), 0)
}

func SetLastUpdateCheck() {
	newValue := time.Now().Unix()
	configUpdated = true
	log.Debugln("config: setting", lastUpdateCheckKey, "to", newValue)
	viper.Set(lastUpdateCheckKey, newValue)
}

func WriteConfig() {
	var err error
	if existingConfig && configUpdated {
		err = viper.WriteConfig()
	} else if !existingConfig {
		err = viper.SafeWriteConfig()
		logger.CheckErr(err)
		err = viper.ReadInConfig() // hack to get ConfigFileUsed working
	}

	if configUpdated || !existingConfig {
		logger.CheckErr(err)
		log.WithFields(rawFields()).Debugln("config: saved to", viper.ConfigFileUsed())
	}

	existingConfig = true
	configUpdated = false
}

func rawFields() log.Fields {
	return log.Fields{
		autoUpdateKey:      viper.GetBool(autoUpdateKey),
		updateIntervalKey:  viper.GetString(updateIntervalKey),
		lastUpdateCheckKey: viper.GetInt64(lastUpdateCheckKey),
	}
}
