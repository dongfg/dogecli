package cmd

import (
	"fmt"

	"github.com/dongfg/dogecli/internal/constants"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of " + constants.CLIName,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version %s\n", constants.CLIName, constants.Version)
	},
}
