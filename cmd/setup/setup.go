package setup

import (
	"fmt"

	"github.com/power-slide/cli/pkg/updater"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "setup",
		Short: "PowerSlide Setup",
		Long:  "Ensures you have everything you need to setup and run PowerSlide.",
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	updater.AutomaticUpdate()
	fmt.Println("You're ready to PowerSlide! 🚀")
}
