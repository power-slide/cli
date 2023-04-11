package versions

import (
	"fmt"

	"github.com/power-slide/cli/pkg/version"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "versions",
		Short: "Get version info",
		Long:  `Get the version info for the cli and the cluster.`,
		Run:   execute,
	}

	localOnly   bool
	latestOnly  bool
	clusterOnly bool
)

func init() {
	Cmd.Flags().BoolVarP(&localOnly, "local-only", "l", false, "Only print cli version")
	Cmd.Flags().BoolVarP(&latestOnly, "latest-only", "r", false, "Only print latest version")
	Cmd.Flags().BoolVarP(&clusterOnly, "cluster-only", "u", false, "Only print cluster version")
	Cmd.MarkFlagsMutuallyExclusive("local-only", "latest-only", "cluster-only")
}

func execute(cmd *cobra.Command, args []string) {
	updater.AutomaticUpdate()

	noFlagSet := !(latestOnly || localOnly || clusterOnly)

	if noFlagSet || localOnly {
		fmt.Printf("PowerSlide CLI version: %s\n", version.Version)
	}

	if noFlagSet || clusterOnly {
		// TODO: get the version number once it's stored somewhere
		fmt.Printf("PowerSlide cluster version: %s\n", "unknown")
	}
}
