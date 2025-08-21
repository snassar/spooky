package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show spooky version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("spooky v0.1.0")
		fmt.Println("Automation and configuration management tool")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
