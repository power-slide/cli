package dev

import (
	"github.com/power-slide/cli/cmd/dev/environment"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "dev",
	Short: "PowerSlide Developer Tools",
	Long:  `Developer tools for building apps on PowerSlide.`,
}

func init() {
	Cmd.AddCommand(environment.Cmd)
}
