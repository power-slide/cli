package create

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/power-slide/cli/cmd/util/colors"
)

var (
	optionalCommands = []string{"asdf", "argocd"}
	requiredCommands = []string{"k3d", "kubectl"}
)

func checkForCommands() {
	fmt.Println("Checking for installed dependencies...")
	fmt.Println("\nOptional commands:")
	for _, command := range optionalCommands {
		printCommandStatus(command, isCommandInstalled(command))
	}

	fmt.Printf("\n\n")
	fmt.Println("Required commands:")
	missingCommand := false
	for _, command := range requiredCommands {
		if installed := isCommandInstalled(command); installed {
			printCommandStatus(command, installed)
		} else {
			missingCommand = true
			printCommandStatus(command, installed)
		}
	}

	if missingCommand {
		fmt.Println("\nPlease install the missing required commands listed above.")
		os.Exit(1)
	}
	fmt.Printf("\nAll needed commands found.\n\n")
}

func printCommandStatus(command string, installed bool) {
	if installed {
		fmt.Printf("%s%s%s is %sinstalled%s.\n", colors.TextYellow, command, colors.TextReset, colors.TextGreen, colors.TextReset)
	} else {
		fmt.Printf("%s%s%s is %sNOT installed%s.\n", colors.TextYellow, command, colors.TextReset, colors.TextRed, colors.TextReset)
	}
}

func isCommandInstalled(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
