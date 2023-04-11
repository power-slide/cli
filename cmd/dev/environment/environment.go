package environment

import (
	"github.com/power-slide/cli/cmd/dev/environment/create"
	"github.com/power-slide/cli/cmd/dev/environment/destroy"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "environment",
	Short: "Manage PowerSlide development environments",
	Long:  `Manage PowerSlide development environments.`,
}

func init() {
	Cmd.AddCommand(create.Cmd)
	Cmd.AddCommand(destroy.Cmd)
}
