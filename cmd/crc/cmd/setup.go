package cmd

import (
	"github.com/code-ready/crc/pkg/crc/preflight"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup hypervisor",
	Long:  "setup hypervisor to run the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		runSetup(args)
	},
}

func runSetup(arguments []string) {
	preflight.SetupHost()
}
