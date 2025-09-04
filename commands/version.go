package commands

import (
	"fmt"

	"spooky/internal/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show spooky version information",
	Run: func(cmd *cobra.Command, args []string) {
		versionInfo := version.GetFullVersionInfo()
		fmt.Println(versionInfo.String())
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
