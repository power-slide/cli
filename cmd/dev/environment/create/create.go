package create

import (
	"fmt"
	"time"

	"github.com/power-slide/cli/cmd/util"
	"github.com/power-slide/cli/cmd/util/colors"
	"github.com/power-slide/cli/pkg/updater"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "create [environment-name]",
		Short: "Create a PowerSlide development environment",
		Long:  `Spins up a PowerSlide development environment locally using k3d by default.`,
		Args:  checkArgs,
		Run:   run,
	}

	skipCluster    bool
	skipMonitoring bool
	skipArgoCD     bool
	clusterPort    int
	cmdTimeoutFlag string
	cmdTimeout     time.Duration
)

func init() {
	Cmd.Flags().BoolVar(&skipCluster, "skip-cluster", false, "Skip creating cluster")
	Cmd.Flags().IntVarP(&clusterPort, "cluster-port", "p", 42069, "PowerSlide services port")
	Cmd.Flags().BoolVar(&skipMonitoring, "skip-monitoring", false, "Skip installing monitoring")
	Cmd.Flags().BoolVar(&skipArgoCD, "skip-argocd", false, "Skip installing argocd")
	Cmd.Flags().StringVarP(&cmdTimeoutFlag, "timeout", "t", "2m", "Time to wait for installation steps to complete")
}

func checkArgs(cmd *cobra.Command, args []string) error {
	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		return err
	}

	var err error
	cmdTimeout, err = time.ParseDuration(cmdTimeoutFlag)
	if err != nil {
		return err
	}

	if clusterPort <= 1024 || clusterPort > 65535 {
		return fmt.Errorf("invalid port: %d", clusterPort)
	}

	return nil
}

func run(cmd *cobra.Command, args []string) {
	updater.AutomaticUpdate()
	fmt.Println("Creating a local PowerSlide environment...")
	checkForCommands()
	createCluster(args[0])
	currentContext := util.Kubectl([]string{"config", "current-context"}, "")
	fmt.Printf("Got kubectl context: %s%s%s\n", colors.TextBlue, currentContext, colors.TextReset)
	installMonitoringStack()
	installArgoCD()
	fmt.Printf("\nReady to %s%s%s!\n", colors.TextGreen, "PowerSlide", colors.TextReset)
}
