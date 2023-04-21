package cmd

import (
	"fmt"

	"github.com/power-slide/cli/cmd/dev"
	"github.com/power-slide/cli/cmd/setup"
	"github.com/power-slide/cli/cmd/update"
	"github.com/power-slide/cli/cmd/versions"
	"github.com/power-slide/cli/pkg/config"
	"github.com/power-slide/cli/pkg/logger"
	versionNumber "github.com/power-slide/cli/pkg/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pwrsl",
	Short: "PowerSlide's multi-tool",
	Long: `PowerSlide is an open source PaaS that runs on Kubernetes.
This tool is for admins and developers to manage PowerSlide
clusters and applications.`,
	Version: versionNumber.Version,
}

func Execute() {
	logger.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(logger.Init)
	cobra.OnInitialize(config.Init)
	cobra.OnFinalize(config.WriteConfig)

	configHelpMsg := fmt.Sprintf(
		"config file (default is $HOME/%s.%s)",
		config.DefaultConfigFile,
		config.DefaultConfigFormat,
	)
	rootCmd.PersistentFlags().StringVarP(&config.ConfigFile, "config", "c", "", configHelpMsg)

	rootCmd.AddCommand(dev.Cmd)
	rootCmd.AddCommand(update.Cmd)
	rootCmd.AddCommand(setup.Cmd)
	rootCmd.AddCommand(versions.Cmd)
}
