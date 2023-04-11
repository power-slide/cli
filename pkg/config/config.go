package config

import (
	"os"

	"github.com/power-slide/cli/pkg/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	}
}
