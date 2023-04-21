package destroy

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/power-slide/cli/pkg/updater"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "destroy",
	Short: "Tear down a PowerSlide development environment",
	Long:  "Tears down a PowerSlide development environment locally using k3d by default.",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	updater.AutomaticUpdate()
	clusterCommandArgs := []string{"cluster", "delete", args[0]}
	fmt.Printf("Deleting local k3s cluster... ")
	k3dCommand := exec.Command("k3d", clusterCommandArgs...)

	if out, err := k3dCommand.CombinedOutput(); err != nil {
		fmt.Println("k3d output: ", string(out))
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Done!")
}
